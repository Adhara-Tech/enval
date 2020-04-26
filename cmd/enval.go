package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

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
		fmt.Println(err)
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

	//TODO path should not be provided this way. If it is relative to the file, may be the file should manage that
	toolsStorage := infra.NewDefaultToolsStorage("../../tool-specs")
	toolsStorageAdapter := adapters.NewDefaultStorageAdapter(toolsStorage)
	systemAdapter := adapters.NewDefaultSystemAdapter()
	theChecker := manifestchecker.NewDefaultManifestChecker(toolsStorageAdapter, systemAdapter)

	result, err := theChecker.Check(*manifest, func(msg manifestchecker.Notification) {
		if !msg.IsToolAvailable {
			fmt.Printf("%s %s %s", notFoundSymbol, toolName(msg.Tool), "Command Not Found")
			return
		}

		if !msg.IsVersionValid {
			fmt.Printf("%s %s:\n%s", invalidSymbol, toolName(msg.Tool), renderVersions(msg.Tool, msg.VersionsFound, msg.VersionValidations))
			return
		}

		fmt.Printf("%s %s:\n%s", validSymbol, toolName(msg.Tool), renderVersions(msg.Tool, msg.VersionsFound, msg.VersionValidations))
	})
	if err != nil {
		return err
	}

	if !result.Ok {
		return errors.New(result.Message)
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

func renderVersions(tool model.ManifestTool, fieldVersions map[string]string, fieldOk map[string]bool) string {

	var buffer bytes.Buffer
	for fieldName, versionConstraint := range tool.Checks {
		ok := fieldOk[fieldName]
		version := fieldVersions[fieldName]
		symbol := validSymbol

		if !ok {
			symbol = invalidSymbol
		}

		buffer.WriteString(fmt.Sprintf("    %s %s(%s): %s\n", symbol, fieldName, versionConstraint, version))
	}

	return buffer.String()
}
