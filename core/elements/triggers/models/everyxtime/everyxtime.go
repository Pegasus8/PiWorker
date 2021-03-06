package everyxtime

import (
	"fmt"
	"time"

	"github.com/Pegasus8/piworker/core/data"
	"github.com/Pegasus8/piworker/core/elements/triggers/shared"
	"github.com/Pegasus8/piworker/core/types"
)

const triggerID = "T4"

var triggerArgs = []shared.Arg{
	{
		ID:   triggerID + "-1",
		Name: "Time",
		Description: "Time of repetition. Format must be '1h10m10s' where" +
			" 'h' = hours, 'm' = minutes and 's' = seconds.",
		ContentType: types.Text,
	},
}

// EveryXTime - Trigger
var EveryXTime = shared.Trigger{
	ID:          triggerID,
	Name:        "Every X Time",
	Description: "",
	Run:         trigger,
	Args:        triggerArgs,
}

var nextExecution = make(map[string]time.Time)

func trigger(args *[]data.UserArg, parentTaskID string) (result bool, err error) {
	if len(*args) != len(triggerArgs) {
		return false, fmt.Errorf("%d arguments were expected and %d were obtained", len(triggerArgs), len(*args))
	}

	// Time
	var timeToWait time.Duration

	for i, arg := range *args {
		if arg.Content == "" {
			return false, fmt.Errorf("argument %d (ID: %s) is empty", i, arg.ID)
		}

		switch arg.ID {
		case triggerArgs[0].ID:
			{
				timeToWait, err = time.ParseDuration(arg.Content)
				if err != nil {
					return false, err
				}
			}
		default:
			return false, shared.ErrUnrecognizedArgID
		}
	}

	// First execution
	if _, exists := nextExecution[parentTaskID]; !exists {
		nextExecution[parentTaskID] = time.Now().Add(timeToWait)

		return false, nil
	}

	if nextExecution[parentTaskID].Unix() <= time.Now().Unix() {
		nextExecution[parentTaskID] = time.Now().Add(timeToWait)

		return true, nil
	}

	return false, nil
}
