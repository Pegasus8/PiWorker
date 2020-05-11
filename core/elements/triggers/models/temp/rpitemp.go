package temp

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Pegasus8/piworker/core/data"
	"github.com/Pegasus8/piworker/core/elements/triggers/shared"
	"github.com/Pegasus8/piworker/core/types"

	"github.com/shirou/gopsutil/host"
)

const triggerID = "T2"

var triggerArgs = []shared.Arg{
	{
		ID:   triggerID + "-1",
		Name: "Expected Temperature",
		Description: "The expected temperature of the host. Must be in" +
			" float format and without the 'ºC'. Example: 55.1.",
		ContentType: types.Float,
	},
}

// RaspberryTemperature - Trigger
var RaspberryTemperature = shared.Trigger{
	ID:          triggerID,
	Name:        "Raspberry's Temperature",
	Description: "If the temperature of the host equals or exceeds the given number, the trigger will be activated.",
	Run:         trigger,
	Args:        triggerArgs,
}

func trigger(args *[]data.UserArg, parentTaskID string) (result bool, err error) {
	if len(*args) != len(triggerArgs) {
		return false, fmt.Errorf("%d arguments were expected and %d were obtained", len(triggerArgs), len(*args))
	}

	// Expected temperature received
	var expectedTemp float64

	for i, arg := range *args {
		if arg.Content == "" {
			return false, fmt.Errorf("argument %d (ID: %s) is empty", i, arg.ID)
		}

		switch arg.ID {
		// Temperature arg
		case triggerArgs[0].ID:
			{
				expectedTemp, err = strconv.ParseFloat(arg.Content, 64)
				if err != nil {
					return false, err
				}
			}

		default:
			return false, shared.ErrUnrecognizedArgID
		}
	}

	st, err := host.SensorsTemperatures()
	if err != nil {
		return false, err
	}

	temperature := func() float64 {
		for _, t := range st {
			if t.SensorKey == "coretemp_packageid0_input" {
				return t.Temperature
			}
		}
		return 0.0
	}()

	if temperature == 0.0 {
		return false, errors.New("SensorKey incompatible with host")
	}

	if temperature >= expectedTemp {
		return true, nil
	}

	return false, nil
}
