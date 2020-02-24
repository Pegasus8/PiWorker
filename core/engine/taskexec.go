package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Pegasus8/piworker/core/data"
	actionsModel "github.com/Pegasus8/piworker/core/elements/actions"
	actionsList "github.com/Pegasus8/piworker/core/elements/actions/models"
	triggersList "github.com/Pegasus8/piworker/core/elements/triggers/models"
	"github.com/Pegasus8/piworker/core/types"
	"github.com/Pegasus8/piworker/core/uservariables"
	"github.com/rs/zerolog/log"
)

func runTaskLoop(taskID string, taskChannel chan data.UserTask) {
	log.Info().Str("taskID", taskID).Msg("Loop started")
	for {
		// Receive the renewed data for the task in question, if there is not data
		// just keep waiting for it.
		taskReceived := <-taskChannel

		triggered, err := runTrigger(taskReceived.TaskInfo.Trigger, taskReceived.TaskInfo.ID)
		if err != nil {
			log.Error().
				Err(err).
				Str("taskID", taskReceived.TaskInfo.ID).
				Msg("Error while trying to run the trigger of the task, stopping the task execution...")
			break
		}
		if triggered {
			if wasRecentlyExecuted(taskReceived.TaskInfo.ID) {
				log.Debug().
					Str("taskID", taskReceived.TaskInfo.ID).
					Msg("The task was recently executed, the trigger stills active. Skipping it...")
				goto skipTaskExecution
			}

			log.Info().
				Str("taskID", taskReceived.TaskInfo.ID).
				Str("triggerID", taskReceived.TaskInfo.Trigger.ID).
				Msg("[%s] Trigger with the ID '%s' activated, running actions...")
			runActions(&taskReceived)

			err = setAsRecentlyExecuted(taskReceived.TaskInfo.ID)
			if err != nil {
				log.Fatal().
					Err(err).
					Str("taskID", taskReceived.TaskInfo.ID).
					Msg("Error when trying to set a task as recently executed")
				break
			}

		skipTaskExecution:
			// Skip the execution of the task but not skip the entire iteration
			// in case of have to do something else with the task.
		} else {
			if wasRecentlyExecuted(taskReceived.TaskInfo.ID) {
				err = setAsReadyToExecuteAgain(taskReceived.TaskInfo.ID)
				if err != nil {
					log.Fatal().
						Err(err).
						Str("taskID", taskReceived.TaskInfo.ID).
						Msg("Error when trying to set a task as ready to execute again")
					break
				}
			}
		}
	}
}

func runTrigger(trigger data.UserTrigger, parentTaskID string) (bool, error) {
	for _, pwTrigger := range triggersList.TRIGGERS {
		if trigger.ID == pwTrigger.ID {
			for _, arg := range trigger.Args {
				// Check if the arg contains a user global variable
				err := searchAndReplaceVariable(&arg, parentTaskID)
				if err != nil {
					return false, err
				}
			}
			result, err := pwTrigger.Run(&trigger.Args, parentTaskID)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
			return false, nil
		}
	}

	return false, fmt.Errorf("The trigger with the ID '%s' cannot be found", trigger.ID)
}

func runActions(task *data.UserTask) {
	log.Info().Str("taskID", task.TaskInfo.ID).Msg("Running actions...")
	startTime := time.Now()

	userActions := &task.TaskInfo.Actions
	previousState := task.TaskInfo.State

	log.Printf("[%s] Changing task state to '%s'\n", task.TaskInfo.ID, data.StateTaskOnExecution)
	// Set task state to on-execution
	err := data.UpdateTaskState(task.TaskInfo.ID, data.StateTaskOnExecution)
	if err != nil {
		log.Fatal().
			Str("taskID", task.TaskInfo.ID).
			Msgf("Error when trying to update the task state to '%s'", data.StateTaskOnExecution)
	}

	var chainedResult *actionsModel.ChainedResult
	var orderN int8 = 0
	for range *userActions {

		for _, userAction := range *userActions {
			if userAction.Order == orderN {

				// Run the action
				for _, action := range actionsList.ACTIONS {
					if userAction.ID == action.ID {
						log.Info().
							Str("taskID", task.TaskInfo.ID).
							Str("actionID", userAction.ID).
							Bool("chained", userAction.Chained).
							Int8("actionOrder", userAction.Order).
							Str("previousResultType", string(chainedResult.ResultType)).
							Str("previousResultContent", chainedResult.Result).
							Msg("Running action")

						for _, arg := range userAction.Args {
							err := searchAndReplaceVariable(&arg, task.TaskInfo.ID)
							if err != nil {
								log.Error().
									Str("taskID", task.TaskInfo.ID).
									Str("actionID", userAction.ID).
									Str("argID", arg.ID).
									Err(err).
									Int8("actionOrder", orderN).
									Msg("Error when searching for a variable on the argument")
								return
							}
						}

						ua, err := replaceArgByCR(chainedResult, &userAction)
						if err != nil {
							log.Error().
								Str("taskID", task.TaskInfo.ID).
								Str("actionID", userAction.ID).
								Err(err).
								Int8("actionOrder", userAction.Order).
								Msg("Error when trying to replace an argument for a variable")
							return
						}
						userAction = *ua

						result, chr, err := action.Run(chainedResult, &userAction, task.TaskInfo.ID)
						// Set the returned chr (chained result) to our main instance of the ChainedResult struct (`chainedResult`).
						// This will be given to the next action (if exists).
						chainedResult = chr
						if err != nil {
							log.Error().
								Str("taskID", task.TaskInfo.ID).
								Str("actionID", userAction.ID).
								Err(err).
								Int8("actionOrder", userAction.Order).
								Msg("Error when running the action")
							return
						}
						if result {
							log.Info().
								Str("taskID", task.TaskInfo.ID).
								Str("actionID", userAction.ID).
								Int8("actionOrder", userAction.Order).
								Msg("Action finished correctly")
						} else {
							log.Printf("[%s] Action in order %d wasn't executed correctly. Aborting task for prevention of future errors...\n",
								task.TaskInfo.ID, userAction.Order)
							log.Warn().
								Str("taskID", task.TaskInfo.ID).
								Str("actionID", userAction.ID).
								Int8("actionOrder", userAction.Order).
								Msg("Action wasn't executed correctly. Aborting task for prevention of future errors...")
							return
						}

						// No need to keep iterating
						break
					}
				}

				orderN++
				break
			}
		}

	}

	// Needed read the actual task state
	updatedData, err := data.ReadData()
	if err != nil {
		log.Fatal().
			Str("taskID", task.TaskInfo.ID).
			Err(err).
			Msg("Error when trying to read the current status of the task")
	}
	updatedTask, _, err := updatedData.GetTaskByID(task.TaskInfo.ID)
	if err != nil {
		log.Fatal().
			Str("taskID", task.TaskInfo.ID).
			Err(err).
			Msg("Error when trying to get the task by its ID")
	}
	lastState := updatedTask.TaskInfo.State
	// If the state has no changes, return to the original state
	if lastState == data.StateTaskOnExecution {
		err = data.UpdateTaskState(task.TaskInfo.ID, previousState)
		if err != nil {
			log.Fatal().
				Str("taskID", task.TaskInfo.ID).
				Str("previousState", string(previousState)).
				Err(err).
				Msg("Error when trying to update the task's state")
		}
	}
	executionTime := time.Since(startTime).String()
	log.Info().
		Str("taskID", task.TaskInfo.ID).
		Str("executionTime", executionTime).
		Msg("Actions executed")
}

func checkForAnUpdate(updateChannel chan bool) {
	dataPath := filepath.Join(data.DataPath, data.Filename)
	var oldModTime time.Time
	var newModTime time.Time
	for range time.Tick(time.Millisecond * 300) {
		fileInfo, err := os.Stat(dataPath)
		if err != nil {
			log.Error().Err(err).Msg("Error when trying to get the file's info of the user data file")
			continue
		}
		// First run
		if oldModTime.IsZero() {
			log.Info().Msg("First run of the data file watchdog, setting variable of comparison")
			oldModTime = fileInfo.ModTime()
		}
		newModTime = fileInfo.ModTime()
		if oldModTime != newModTime {
			log.Info().Msg("Change detected on the data file, sending the signal...")
			// Send the signal
			updateChannel <- true
			// Update the variable
			oldModTime = newModTime
		}
	}
}

func setAsRecentlyExecuted(ID string) error {
	dir, err := ioutil.TempDir(TempDir, "")
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(filepath.Join(dir, ID), "")
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func wasRecentlyExecuted(ID string) bool {
	_, err := os.Stat(filepath.Join(TempDir, ID))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else if os.IsExist(err) {
			return true
		}
		log.Fatal().Err(err).Str("taskID", ID).Msg("Error when trying to get the the execution file's info")
	}

	return true
}

func setAsReadyToExecuteAgain(ID string) error {
	path := filepath.Join(TempDir, ID)
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func searchAndReplaceVariable(arg *data.UserArg, parentTaskID string) error {
	// Check if the arg contains a user global variable
	if uservariables.ContainGlobalVariable(&arg.Content) {
		// If yes, then get the name of the variable by using regex
		varName := uservariables.GetGlobalVariableName(arg.Content)
		// Get the variable from the name
		globalVar, err := uservariables.GetGlobalVariable(varName)
		if err != nil {
			log.Error().Err(err).Str("taskID", parentTaskID).Str("varName", varName).Msg("Error when trying to read the user global variable")
			return err
		}
		globalVar.RLock()
		// If all it's ok, replace the content of the argument (wich is the variable name basically)
		// with the content of the desired user global variable.
		arg.Content = globalVar.Content
		globalVar.RUnlock()
		// If the arg not contains a user global variable, then check if contains a user local variable instead.
	} else if uservariables.ContainLocalVariable(&arg.Content) {
		// If yes, then get the name of the variable by using regex
		varName := uservariables.GetLocalVariableName(arg.Content)
		// Get the variable from the name
		localVariable, err := uservariables.GetLocalVariable(varName, parentTaskID)
		if err != nil {
			log.Error().Err(err).Str("taskID", parentTaskID).Str("varName", varName).Msg("Error when trying to read the user local variable")
			return err
		}
		localVariable.RLock()
		// If all it's ok, replace the content of the argument (which is the variable name basically)
		// with the content of the desired user local variable.
		arg.Content = localVariable.Content
		localVariable.RUnlock()
	}

	return nil
}

func replaceArgByCR(chainedResult *actionsModel.ChainedResult, userAction *data.UserAction) (*data.UserAction, error) {
	if userAction.Order == 0 {
		// Prevent the usage of ChainedResult because there are no previous actions.
		userAction.Chained = false
	}
	if userAction.Chained {
		if chainedResult.Result == "" {
			return nil, actionsList.ErrEmptyChainedResult
		}

		for _, userArg := range userAction.Args {
			if userArg.ID == userAction.ArgumentToReplaceByCR {
				userArgType, err := getUserArgType(userAction.ID, userArg.ID)
				if err != nil {
					return nil, err
				}
				if chainedResult.ResultType != userArgType && userArgType != types.Any {
					return nil, fmt.Errorf("Can't replace the arg with the ID '%s' of type '%s' with the previous ChainedResult of type '%s'", userArg.ID, userArgType, chainedResult.ResultType)

				}
				// If all is ok, replace the content
				userArg.Content = chainedResult.Result
			}
		}
	}

	return userAction, nil
}

func getUserArgType(userActionID string, userArgID string) (types.PWType, error) {
	var actionFound bool

	for _, action := range actionsList.ACTIONS {
		if action.ID == userActionID {
			actionFound = true
			for _, arg := range action.Args {
				if arg.ID == userArgID {
					return arg.ContentType, nil
				}
			}
		}
	}

	var err error
	if actionFound {
		err = fmt.Errorf("Unrecognized argument ID '%s' of the action '%s'", userArgID, userActionID)
	} else {
		err = fmt.Errorf("Unrecognized action ID '%s'", userActionID)
	}

	return types.Any, err
}