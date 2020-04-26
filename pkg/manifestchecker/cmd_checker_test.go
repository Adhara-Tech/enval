package manifestchecker_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/infra"
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
	"github.com/Adhara-Tech/enval/pkg/model"
)

type testSystemAdapter struct {
	outputFilePath string
}

func (t *testSystemAdapter) NextOutput(outputFilePath string) {
	t.outputFilePath = outputFilePath
}

func (t testSystemAdapter) CheckCommandAvailable(command string) (bool, error) {
	return true, nil
}

func (t *testSystemAdapter) GetCommandVersionOutput(commandName string, params []string) (string, error) {
	if t.outputFilePath == "" {
		return "", errors.New("file output path is empty")
	}
	data, err := ioutil.ReadFile(t.outputFilePath)
	if err != nil {
		return "", err
	}

	t.outputFilePath = ""

	return string(data), nil
}

var _ manifestchecker.SystemAdapter = (*testSystemAdapter)(nil)

func TestDefaultManifestChecker_Check(t *testing.T) {
	toolsStorage := infra.NewDefaultToolsStorage("../../tool-specs")
	toolsStorageAdapter := adapters.NewDefaultStorageAdapter(toolsStorage)
	systemAdapter := &testSystemAdapter{}
	checker := manifestchecker.NewDefaultManifestChecker(toolsStorageAdapter, systemAdapter)

	systemAdapter.NextOutput("testdata/tools-version-output/go_1.14.output.txt")

	manifest := model.Manifest{
		Name:  "go-tool-test",
		Tools: []model.ManifestTool{{Name: "go", Checks: map[string]string{"version": ">= 1.10"}}},
	}
	err := checker.Check(manifest, func(notification manifestchecker.Notification) {
		fmt.Println(notification)
	})

	require.Nil(t, err)
}
