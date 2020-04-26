package adapters

import (
	"os/exec"

	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
)

type DefaultSystemAdapter struct {
}

func NewDefaultSystemAdapter() *DefaultSystemAdapter {
	return &DefaultSystemAdapter{}
}

func (systemAdapter DefaultSystemAdapter) CheckCommandAvailable(command string) (bool, error) {
	_, err := exec.LookPath(command)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (systemAdapter DefaultSystemAdapter) GetCommandVersionOutput(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)

	versionString, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(versionString), nil
}

var _ manifestchecker.SystemAdapter = (*DefaultSystemAdapter)(nil)
