package adapters

import (
	"os"
	"os/exec"

	"github.com/Adhara-Tech/enval/pkg/exerrors"

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
		if execError, ok := err.(*exec.Error); ok {
			if execError.Err == exec.ErrNotFound {
				return false, nil
			}
		}
		return false, exerrors.Wrap(err, exerrors.InternalEnvalErrorKind)
	}

	return true, nil
}

func (systemAdapter DefaultSystemAdapter) ExecuteCommand(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)

	versionString, err := cmd.CombinedOutput()
	if err != nil {
		return "", exerrors.Wrap(err, exerrors.InternalEnvalErrorKind)
	}

	return string(versionString), nil
}

func (systemAdapter DefaultSystemAdapter) CheckDirExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, exerrors.Wrap(err, exerrors.InternalEnvalErrorKind)
	}

	if !fileInfo.IsDir() {
		return false, nil
	}

	return true, nil
}

var _ manifestchecker.SystemAdapter = (*DefaultSystemAdapter)(nil)
