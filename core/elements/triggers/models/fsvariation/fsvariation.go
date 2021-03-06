package fsvariation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Pegasus8/piworker/core/data"
	"github.com/Pegasus8/piworker/core/elements/triggers/shared"
	"github.com/Pegasus8/piworker/core/types"
)

const triggerID = "T3"

var triggerArgs = []shared.Arg{
	{
		ID:          triggerID + "-1",
		Name:        "Path of the Objective File",
		Description: "Must be on the format 'path/of/the/file.txt'.",
		ContentType: types.Path,
	},
}

// VariationOfFileSize - Trigger
var VariationOfFileSize = shared.Trigger{
	ID:          triggerID,
	Name:        "Variation of a File's Size",
	Description: "",
	Run:         trigger,
	Args:        triggerArgs,
}

var previousFileSize = make(map[string]int64)

func trigger(args *[]data.UserArg, parentTaskID string) (result bool, err error) {
	if len(*args) != len(triggerArgs) {
		return false, fmt.Errorf("%d arguments were expected and %d were obtained", len(triggerArgs), len(*args))
	}

	// Filepath
	var filePath string

	for i, arg := range *args {
		if arg.Content == "" {
			return false, fmt.Errorf("argument %d (ID: %s) is empty", i, arg.ID)
		}

		switch arg.ID {
		case triggerArgs[0].ID:
			filePath = filepath.Clean(arg.Content)
		default:
			return false, shared.ErrUnrecognizedArgID
		}
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	// First execution
	if _, exists := previousFileSize[parentTaskID]; !exists {
		previousFileSize[parentTaskID] = info.Size()
		return false, nil
	}

	if info.Size() != previousFileSize[parentTaskID] {
		// Update the stored size
		previousFileSize[parentTaskID] = info.Size()
		return true, nil
	}

	return false, nil
}
