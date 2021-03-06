package time

import (
	"fmt"
	"time"

	"github.com/Pegasus8/piworker/core/data"
	"github.com/Pegasus8/piworker/core/elements/triggers/shared"
	"github.com/Pegasus8/piworker/core/types"
)

const triggerID = "T1"

var triggerArgs = []shared.Arg{
	{
		ID:   triggerID + "-1",
		Name: "Date",
		Description: "The date to launch the trigger. The format used is YYYY-MM-dd." +
			" Example: 2019-11-15.",
		ContentType: types.Date,
	},
	{
		ID:   triggerID + "-2",
		Name: "Hour",
		Description: "The hour to launch the  trigger. The format used is HH:mm." +
			" Example: 13:45",
		ContentType: types.Time,
	},
}

// ByTime - Trigger
var ByTime = shared.Trigger{
	ID:          triggerID,
	Name:        "By Time",
	Description: "",
	Run:         trigger,
	Args:        triggerArgs,
}

func trigger(args *[]data.UserArg, parentTaskID string) (result bool, err error) {
	if len(*args) != len(triggerArgs) {
		return false, fmt.Errorf("%d arguments were expected and %d were obtained", len(triggerArgs), len(*args))
	}

	// Contains the time and date received from the arguments.
	var t time.Time
	var date, hour string

	for i, arg := range *args {
		if arg.Content == "" {
			return false, fmt.Errorf("argument %d (ID: %s) is empty", i, arg.ID)
		}

		switch arg.ID {
		case triggerArgs[0].ID:
			date = arg.Content
		case triggerArgs[1].ID:
			hour = arg.Content
		default:
			return false, shared.ErrUnrecognizedArgID
		}
	}

	t, err = time.Parse("2006-01-02 15:04", date+" "+hour)
	if err != nil {
		return false, err
	}

	if time.Now().Format("2006-01-02 15:04") == t.Format("2006-01-02 15:04") {
		return true, nil
	}

	return false, nil
}
