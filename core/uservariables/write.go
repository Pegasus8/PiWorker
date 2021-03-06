package uservariables

import (
	"encoding/json"
	"os"

	"github.com/Pegasus8/piworker/utilities/files"
	"github.com/rs/zerolog/log"
)

// Init initializes the directory where the user variables will be stored (if not exists).
func Init() {
	// Create data path if not exists
	err := os.MkdirAll(UserVariablesPath, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize the directory to store the user variables")
	}
}

// WriteToFile writes the current content of the LocalVariable to the corresponding file.
func (localVar *LocalVariable) WriteToFile() error {
	localVar.Lock()
	filename := localVar.Name + "-" + localVar.ParentTaskID

	byteData, err := json.MarshalIndent(localVar, "", "   ")
	if err != nil {
		localVar.Unlock()
		return err
	}
	localVar.Unlock()

	_, err = files.WriteFile(UserVariablesPath, filename, byteData)
	if err != nil {
		return err
	}

	return nil
}

// WriteToFile writes the current content of the GlobalVariable to the corresponding file.
func (globalVar *GlobalVariable) WriteToFile() error {
	globalVar.Lock()
	byteData, err := json.MarshalIndent(globalVar, "", "   ")
	if err != nil {
		globalVar.Unlock()
		return err
	}
	globalVar.Unlock()

	_, err = files.WriteFile(UserVariablesPath, globalVar.Name, byteData)
	if err != nil {
		return err
	}

	return nil
}
