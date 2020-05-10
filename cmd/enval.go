package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Adhara-Tech/enval/pkg/exerrors"

	"github.com/Adhara-Tech/enval/cmd/version"
	"github.com/Adhara-Tech/enval/pkg/adapters"
	"github.com/Adhara-Tech/enval/pkg/config"
	"github.com/Adhara-Tech/enval/pkg/infra"
	"github.com/Adhara-Tech/enval/pkg/manifestchecker"
	"github.com/Adhara-Tech/enval/pkg/model"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name = "enval"
)

var cmd = &cobra.Command{
	Use:  "enval",
	Long: name,
	RunE: executeCmd,
}

const (
	manifestFlag = "manifest"
)

func main() {

	cmd.Flags().String(manifestFlag, "", "path to the manifest file")

	err := viper.BindPFlag(manifestFlag, cmd.Flags().Lookup(manifestFlag))
	if err != nil {
		panic(err)
	}

	cmd.AddCommand(version.Command())

	viper.AutomaticEnv()
	if err := cmd.Execute(); err != nil {
		fmt.Println(exerrors.ErrorStack(err))
		os.Exit(1)
	}
}

func executeCmd(_ *cobra.Command, _ []string) error {

	//fmt.Println(version, commitHash, buildTime, branch)

	var manifest *model.Manifest
	var err error

	if viper.IsSet(manifestFlag) {
		manifestPath := viper.GetString(manifestFlag)
		manifest, err = config.ReadManifestFrom(manifestPath)
	} else {
		manifest, err = config.ReadManifest()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	toolsStorage := infra.NewDefaultToolsStorage()
	if manifest.CustomSpecs != "" {
		toolsStorage = toolsStorage.WithCustomSpecs(manifest.CustomSpecs)
	}
	toolsStorageAdapter := adapters.NewDefaultStorageAdapter(toolsStorage)
	systemAdapter := adapters.NewDefaultSystemAdapter()
	versionValidators := map[string]manifestchecker.FieldVersionValidator{
		"semver": manifestchecker.SemverFieldVersionValidator{},
	}
	fieldVersionValidatorManager := manifestchecker.NewFieldVersionValidatorManager(versionValidators)
	versionCheckerManager := manifestchecker.NewVersionCheckerManager(fieldVersionValidatorManager)
	toolsManager := manifestchecker.NewToolsManager(toolsStorageAdapter, systemAdapter, versionCheckerManager)

	_, err = toolsManager.ValidateManifestAndNotify(*manifest, cmdNotifier)
	if err != nil {
		return err
	}

	return nil
}

var validSymbol = color.GreenString("✔")
var invalidSymbol = color.RedString("!")
var notFoundSymbol = color.RedString("∅")

func toolName(tool model.ManifestTool) string {
	if tool.Flavor != nil {
		return fmt.Sprintf("%s(%s)", tool.Name, *tool.Flavor)
	}
	return tool.Name
}

func renderVersions(tool model.ManifestTool, fieldVersions map[string]manifestchecker.FieldValidationResult) string {

	var buffer bytes.Buffer
	for fieldName, versionConstraint := range tool.Checks {
		fieldValidationResult := fieldVersions[fieldName]
		ok := fieldValidationResult.IsValid
		version := fieldValidationResult.ValueFound
		symbol := validSymbol

		if !ok {
			symbol = invalidSymbol
		}

		buffer.WriteString(fmt.Sprintf("    %s %s(%s): %s\n", symbol, fieldName, versionConstraint, version))
	}

	return buffer.String()
}

func cmdNotifier(validationResultArr []manifestchecker.ToolValidationResult) {
	for _, toolValidation := range validationResultArr {
		if !toolValidation.IsToolAvailable {
			fmt.Printf("%s %s: %s\n", notFoundSymbol, toolName(toolValidation.Tool), "Command Not Found")
			return
		}

		if !toolValidation.IsVersionValid {
			fmt.Printf("%s %s: %s\n%s", invalidSymbol, toolName(toolValidation.Tool), toolValidation.ResultDescription, renderVersions(toolValidation.Tool, toolValidation.FieldValidations))
			return
		}

		fmt.Printf("%s %s:\n%s", validSymbol, toolName(toolValidation.Tool), renderVersions(toolValidation.Tool, toolValidation.FieldValidations))
	}
}
