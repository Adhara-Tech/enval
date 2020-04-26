package manifestchecker_test

import (
	"errors"
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

func TestDefaultManifestChecker_CheckValidVersions(t *testing.T) {
	//TODO maybe externalize testdata to a file. This can be huge at some point
	testData := []struct {
		TestName              string
		Manifest              model.Manifest
		VersionOutputFilePath string
	}{
		{
			TestName:              "Go language",
			Manifest:              ManifestFrom("go", map[string]string{"version": ">= 1.14"}),
			VersionOutputFilePath: "testdata/tools-version-output/go_1.14.output.txt",
		},
		{
			TestName:              "golangci-lint go metalinter",
			Manifest:              ManifestFrom("golangci-lint", map[string]string{"version": ">= 1.24.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/golangci-lint_1.24.output.txt",
		},
		{
			TestName:              "Java openjdk ga 11",
			Manifest:              ManifestFromWithFlavor("java", "openjdk", map[string]string{"version": ">= 11.0.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/java_openjdk_11.output.txt",
		},
		{
			TestName:              "Java openjdk ga 14",
			Manifest:              ManifestFromWithFlavor("java", "openjdk", map[string]string{"version": ">= 14.0.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/java_openjdk_14.output.txt",
		},
	}

	toolsStorage := infra.NewDefaultToolsStorage("../../tool-specs")
	toolsStorageAdapter := adapters.NewDefaultStorageAdapter(toolsStorage)
	systemAdapter := &testSystemAdapter{}
	checker := manifestchecker.NewDefaultManifestChecker(toolsStorageAdapter, systemAdapter)

	for index := range testData {
		currentTestData := testData[index]
		t.Run(currentTestData.TestName, func(t *testing.T) {
			systemAdapter.NextOutput(currentTestData.VersionOutputFilePath)

			manifest := currentTestData.Manifest
			result, err := checker.Check(manifest, func(notification manifestchecker.Notification) {
				require.True(t, notification.IsToolAvailable)
				require.True(t, notification.IsVersionValid)
			})

			require.Nil(t, err)
			require.NotNil(t, result)
			require.True(t, result.Ok)
		})
	}

}

func TestDefaultManifestChecker_CommandDoesNotExist(t *testing.T) {
	t.Skip("missing test")
}

func TestDefaultManifestChecker_CommandNotConfigured(t *testing.T) {
	t.Skip("missing test")
}

func TestDefaultManifestChecker_VersionOutputDoesNotMatchRegexp(t *testing.T) {
	t.Skip("missing test")
}

func TestDefaultManifestChecker_FlavorNotSetInManifest(t *testing.T) {
	t.Skip("missing test")
}
