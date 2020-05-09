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
		{
			TestName:              "Node 12.13.1",
			Manifest:              manifestchecker.ManifestFrom("node", map[string]string{"version": ">= 12.13.1"}),
			VersionOutputFilePath: "testdata/tools-version-output/node_v12.13.1.output.txt",
		},
		{
			TestName:              "npm 6.12.1",
			Manifest:              manifestchecker.ManifestFrom("npm", map[string]string{"version": ">= 6.12.1"}),
			VersionOutputFilePath: "testdata/tools-version-output/npm_6.12.1.output.txt",
		},
		{
			TestName:              "truffle 5.1.13",
			Manifest:              manifestchecker.ManifestFrom("truffle", map[string]string{"version": ">= 5.1.13", "solidity": ">= 0.5.16"}),
			VersionOutputFilePath: "testdata/tools-version-output/truffle_5.1.13.output.txt",
		},
		{
			TestName:              "gotestsum 0.4.0",
			Manifest:              manifestchecker.ManifestFrom("gotestsum", map[string]string{"version": "= 0.4.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/gotestsum_0.4.0.output.txt",
		},
		{
			TestName:              "terraform 0.12.24",
			Manifest:              manifestchecker.ManifestFrom("terraform", map[string]string{"version": "= 0.12.24"}),
			VersionOutputFilePath: "testdata/tools-version-output/terraform_v0.12.24.output.txt",
		},
		{
			TestName:              "openapi-generator 4.3.0",
			Manifest:              manifestchecker.ManifestFrom("openapi-generator", map[string]string{"version": "= 4.3.0"}),
			VersionOutputFilePath: "testdata/tools-version-output/openapi-generator_4.3.0.output.txt",
		},
		{
			TestName:              "swagger-cli 4.0.2",
			Manifest:              manifestchecker.ManifestFrom("swagger-cli", map[string]string{"version": "= 4.0.2"}),
			VersionOutputFilePath: "testdata/tools-version-output/swagger-cli.4.0.2.output.txt",
		},
		{
			TestName:              "ruby 2.6.3",
			Manifest:              manifestchecker.ManifestFrom("ruby", map[string]string{"version": "= 2.6.3"}),
			VersionOutputFilePath: "testdata/tools-version-output/ruby_2.6.3.output.txt",
		},
	}

	toolsStorage := infra.NewDefaultToolsStorage()
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
