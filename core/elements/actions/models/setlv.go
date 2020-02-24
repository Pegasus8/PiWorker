package models

import (
	"errors"
	"github.com/Pegasus8/piworker/core/data"
	"github.com/Pegasus8/piworker/core/elements/actions"
	"github.com/Pegasus8/piworker/core/types"
	"github.com/Pegasus8/piworker/core/uservariables"
	"strings"
)

const (
	// Action
	setLocalVariableID = "A5"

	// Args
	variableNameSetLocalVariableID    = setLocalVariableID + "-1"
	variableContentSetLocalVariableID = setLocalVariableID + "-2"
)

// SetLocalVariable - Action
var SetLocalVariable = actions.Action{
	ID:          setLocalVariableID,
	Name:        "Set Local Variable",
	Description: "Sets the content of a local variable. If the variable does not exist, it will be created.",
	Run:         setLocalVariableAction,
	Args: []actions.Arg{
		actions.Arg{
			ID:   variableNameSetLocalVariableID,
			Name: "Name",
			Description: "The name of the variable. Must be lowercase, without spaces or special characters. " +
				"The unique special character allowed is the underscore ('_'). Remind that local variables are only " +
				"valid on the task where are created. If you want share a variable between tasks use a global variable" +
				" instead. Example of variable: some_random_local_var",
			ContentType: types.Text,
		},
		actions.Arg{
			ID:   variableContentSetLocalVariableID,
			Name: "Variable content",
			Description: "The content of the variable. Optionally can be: a result of a previous action, " +
				"another variable or static content (setted by you).",
			ContentType: types.Any,
		},
	},
	ReturnedChainResultDescription: "The content setted to the variable.",
	ReturnedChainResultType:        types.Any,
}

func setLocalVariableAction(previousResult *actions.ChainedResult, parentAction *data.UserAction, parentTaskID string) (result bool, chainedResult *actions.ChainedResult, err error) {
	var args *[]data.UserArg

	// The name of the variable
	var variableName string
	// The content of the variable
	var variableContent string

	args = &parentAction.Args

	for _, arg := range *args {
		switch arg.ID {
		case variableNameSetGlobalVariableID:
			{
				variableName = strings.TrimSpace(arg.Content)
			}
		case variableContentSetGlobalVariableID:
			variableContent = arg.Content
		default:
			{
				return false, &actions.ChainedResult{}, ErrUnrecognizedArgID
			}
		}
	}

	if variableName == "" || variableContent == "" {
		return false, &actions.ChainedResult{}, errors.New("Error: variableName or variableContent empty")
	}

	variableType := types.GetType(variableContent)

	lv := &uservariables.LocalVariable{
		Name:         variableName,
		Content:      variableContent,
		Type:         variableType,
		ParentTaskID: parentTaskID,
	}
	err = lv.WriteToFile()
	if err != nil {
		return false, &actions.ChainedResult{}, err
	}

	var varExists bool
	var index int
	for i, variable := range *uservariables.LocalVariablesSlice {
		if variable.Name == variableName {
			varExists = true
			index = i
		}
	}
	if varExists {
		// If the variable already exists, replace it.
		(*uservariables.LocalVariablesSlice)[index] = *lv
	} else {
		// If the variable does not exists, add it.
		// We can't append directly to the LVS because the function append doesn't returns a pointer.
		newLVS := append(*uservariables.LocalVariablesSlice, *lv)
		uservariables.LocalVariablesSlice = &newLVS
	}

	return true, &actions.ChainedResult{Result: variableContent, ResultType: variableType}, nil
}