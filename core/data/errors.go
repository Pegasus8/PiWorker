package data

import (
	"errors"
)

// ErrBadTaskName is an error used when a task with specific name is not found.
// in the JSON data file.
var ErrBadTaskName = errors.New("Invalid task name: the task name provided not exists " +
	"in the user database.")

// ErrBadTaskID is an error used when a task with specific ID is not found.
// in the JSON data file.
var ErrBadTaskID = errors.New("Invalid task ID: the task ID provided not exists " +
	"in the user database.")

// ErrNoFilenameAssigned is an error used when the name of the json data file was not setted.
var ErrNoFilenameAssigned = errors.New("No Filename: the filename of the data file was" +
	" not assigned")

// ErrBackupLoopAlreadyActive is the error used when the backup loop was already started
// and is called to start again.
var ErrBackupLoopAlreadyActive = errors.New(
	"Error: the backup loop is already active, new loop aborted",
)