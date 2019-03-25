package data


//					** Data storage **					//

// Filename is the name of the user JSON data file. Must be dynamically assigned
// depending of the size of the file.
var Filename = ""
const (
	// DataPath is the path of the user JSON data file.
	DataPath string = "./data/" 
)


//					** Tasks's States **					//

const (
	// StateTaskCompleted is a variable that can be used in the `State` 
	// field of every task. This state represents a finished task.
	StateTaskCompleted string = "completed"

	// StateTaskOnExecution is a variable that can be used in the `State` 
	// field of every task. This state represents a task currently on execution.
	StateTaskOnExecution string = "on-execution"

	// StateTaskInactive is a variable that can be used in the `State` 
	// field of every task. This state represents a deactivated/inactive task.
	StateTaskInactive string = "inactive"

	// StateTaskActive is a variable that can be used in the `State` 
	// field of every task. This state represents an active task.
	StateTaskActive string = "active"
)
