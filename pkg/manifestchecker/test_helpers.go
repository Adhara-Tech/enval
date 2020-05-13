package manifestchecker

import (
	"errors"
	"io/ioutil"

	"github.com/Adhara-Tech/enval/pkg/model"
)

type TestSystemAdapter struct {
	outputFilePath string
}

func (t *TestSystemAdapter) CheckDirExist(path string) (bool, error) {
	panic("implement me")
}

func (t *TestSystemAdapter) NextOutput(outputFilePath string) {
	t.outputFilePath = outputFilePath
}

func (t *TestSystemAdapter) ClearOutput() {
	t.outputFilePath = ""
}

func (t TestSystemAdapter) CheckCommandAvailable(command string) (bool, error) {
	return true, nil
}

func (t *TestSystemAdapter) ExecuteCommand(commandName string, params []string) (string, error) {
	if t.outputFilePath == "" {
		return "", errors.New("test system adapter file output path is empty")
	}
	data, err := ioutil.ReadFile(t.outputFilePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

var _ SystemAdapter = (*TestSystemAdapter)(nil)

func manifest(toolName string, flavor *string, checks map[string]string) model.Manifest {
	return model.Manifest{
		Name: "demo manifest",
		Tools: []model.ManifestTool{
			{
				Name:   toolName,
				Flavor: flavor,
				Checks: checks,
			},
		},
	}
}

func ManifestFrom(toolName string, checks map[string]string) model.Manifest {
	return manifest(toolName, nil, checks)
}

func ManifestFromWithFlavor(toolName string, flavor string, checks map[string]string) model.Manifest {
	return manifest(toolName, &flavor, checks)
}
