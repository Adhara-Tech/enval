package manifestchecker_test

import (
	"fmt"
	"testing"

	"github.com/Adhara-Tech/enval/pkg/exerrors"

	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/infra"
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
	"github.com/Adhara-Tech/enval/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestToolsManager_Validate(t *testing.T) {
	//TODO maybe externalize testdata to a file. This can be huge at some point
	testData := []struct {
		TestName              string
		Manifest              model.Manifest
		VersionOutputFilePath string
	}{
		{
			TestName:              "Go language",
			Manifest:              manifestchecker.ManifestFrom("go", map[string]string{"version": ">= 1.14"}),
			VersionOutputFilePath: "testdata/tools-version-output/go_1.14.output.txt",
		},
		{
			TestName:              "golangci-lint go metalinter",
			Manifest:              manifestchecker.ManifestFrom("golangci-lint", map[string]string{"version": ">= 1.24.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/golangci-lint_1.24.output.txt",
		},
		{
			TestName:              "Java openjdk ga 11",
			Manifest:              manifestchecker.ManifestFromWithFlavor("java", "openjdk", map[string]string{"version": ">= 11.0.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/java_openjdk_11.output.txt",
		},
		{
			TestName:              "Java openjdk ga 14",
			Manifest:              manifestchecker.ManifestFromWithFlavor("java", "openjdk", map[string]string{"version": ">= 14.0.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/java_openjdk_14.output.txt",
		},
		{
			TestName:              "Java openjdk ea 14",
			Manifest:              manifestchecker.ManifestFromWithFlavor("java", "openjdk", map[string]string{"version": ">= 14"}),
			VersionOutputFilePath: "testdata/tools-version-output/java_openjdk_14-ea.output.txt",
		},
		{
			TestName:              "Java ga 14 without flavor",
			Manifest:              manifestchecker.ManifestFrom("java", map[string]string{"version": ">= 14"}),
			VersionOutputFilePath: "testdata/tools-version-output/java_openjdk_14.output.txt",
		},
	}

	toolsStorage := infra.NewDefaultToolsStorage("../../tool-specs")
	toolsStorageAdapter := adapters.NewDefaultStorageAdapter(toolsStorage)
	systemAdapter := &manifestchecker.TestSystemAdapter{}
	versionValidators := map[string]manifestchecker.FieldVersionValidator{
		"semver": manifestchecker.SemverFieldVersionValidator{},
	}
	fieldVersionValidatorManager := manifestchecker.NewFieldVersionValidatorManager(versionValidators)
	versionCheckerManager := manifestchecker.NewVersionCheckerManager(fieldVersionValidatorManager)

	toolsCheckerManager := manifestchecker.NewToolsManager(toolsStorageAdapter, systemAdapter, versionCheckerManager)

	for index := range testData {
		currentTestData := testData[index]
		t.Run(currentTestData.TestName, func(t *testing.T) {
			systemAdapter.NextOutput(currentTestData.VersionOutputFilePath)

			result, err := toolsCheckerManager.ValidateManifest(currentTestData.Manifest)
			if err != nil {
				fmt.Println(exerrors.ErrorStack(err))
			}
			require.Nil(t, err)
			require.NotNil(t, result)

			for _, validationResult := range result {

				require.True(t, validationResult.IsVersionValid)
			}

		})
	}
}
