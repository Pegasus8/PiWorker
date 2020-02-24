package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// ReadData is a func that returns the user data into structs
func ReadData() (*UserData, error) {
	fullpath := filepath.Join(DataPath, Filename)
	if err := checkFile(fullpath); err != nil {
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()
	log.Info().Str("path", DataPath).Msg("Reading user data...")

	jsonData, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer jsonData.Close()

	byteContent, err := ioutil.ReadAll(jsonData)
	if err != nil {
		return nil, err
	}

	var data UserData
	err = json.Unmarshal(byteContent, &data)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("User data loaded")
	return &data, nil
}

// GetTaskByName is a method of the UserData struct that returns a specific task,
// searching it by it name.
func (data *UserData) GetTaskByName(name string) (taskFound *UserTask, indexPosition int, err error) {
	for index, task := range data.Tasks {
		if task.TaskInfo.Name == name {
			return &data.Tasks[index], index, nil
		}
	}
	return nil, 0, ErrBadTaskName
}

// GetTaskByID is a method of the UserData struct that returns a specific task,
// searching it by it ID.
func (data *UserData) GetTaskByID(ID string) (taskFound *UserTask, indexPosition int, err error) {
	for index, task := range data.Tasks {
		if task.TaskInfo.ID == ID {
			return &data.Tasks[index], index, nil
		}
	}
	return nil, 0, ErrBadTaskID
}

// GetActiveTasks is a method of the UserData struct that returns the tasks
// with the state `active`.
func (data *UserData) GetActiveTasks() (activeTasks *[]UserTask) {
	at := make([]UserTask, 0)
	for _, userTask := range data.Tasks {
		if userTask.TaskInfo.State == StateTaskActive {
			at = append(at, userTask)
		}
	}

	return &at
}

// GetInactiveTasks is a method of the UserData struct that returns the tasks
// with the state `inactive`.
func (data *UserData) GetInactiveTasks() (inactiveTasks *[]UserTask) {
	it := make([]UserTask, 0)
	for _, userTask := range data.Tasks {
		if userTask.TaskInfo.State == StateTaskInactive {
			it = append(it, userTask)
		}
	}

	return &it
}

// GetCompletedTasks is a method of the UserData struct that returns the tasks
// with the state `completed`.
func (data *UserData) GetCompletedTasks() (completedTasks *[]UserTask) {
	ct := make([]UserTask, 0)
	for _, userTask := range data.Tasks {
		if userTask.TaskInfo.State == StateTaskActive {
			ct = append(ct, userTask)
		}
	}

	return &ct
}

// GetOnExecutionTasks is a method of the UserData struct that returns the tasks
// with the state `on-execution`.
func (data *UserData) GetOnExecutionTasks() (onExecutionTasks *[]UserTask) {
	oet := make([]UserTask, 0)
	for _, userTask := range data.Tasks {
		if userTask.TaskInfo.State == StateTaskActive {
			oet = append(oet, userTask)
		}
	}

	return &oet
}