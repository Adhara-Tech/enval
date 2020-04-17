package main

import (
	"Adhara-Tech/check-my-setup/pkg/adapters"
	"Adhara-Tech/check-my-setup/pkg/config"
	"Adhara-Tech/check-my-setup/pkg/infra"
	"Adhara-Tech/check-my-setup/pkg/manifestchecker"
	"Adhara-Tech/check-my-setup/pkg/model"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name = "tchecker"
)

var cmd = &cobra.Command{
	Use:  "tchecker",
	Long: name,
	RunE: executeCmd,
}

const (
	//versionFlag   = "version"
	//buildInfoFlag = "build-info"
	manifestFlag = "manifest"
)

var version = "SNAPSHOT"
var branch = ""
var commitHash = ""
var buildTime = ""

func main() {

	cmd.Flags().String(manifestFlag, "", "path to the manifest file")

	err := viper.BindPFlag(manifestFlag, cmd.Flags().Lookup(manifestFlag))
	if err != nil {
		panic(err)
	}

	//cmd.PersistentFlags().String(config.ContextFlag, "", "context")
	//
	//err := viper.BindPFlag(config.ContextFlag, cmd.PersistentFlags().Lookup(config.ContextFlag))
	//if err != nil {
	//	panic(err)
	//}

	//cmd.Flags().Bool(versionFlag, true, "gets ethgwctl version")

	//err = viper.BindPFlag(versionFlag, cmd.Flags().Lookup(versionFlag))
	//if err != nil {
	//	panic(err)
	//}
	//
	//cmd.Flags().Bool(buildInfoFlag, true, "prints build info")
	//
	//err = viper.BindPFlag(buildInfoFlag, cmd.Flags().Lookup(buildInfoFlag))
	//if err != nil {
	//	panic(err)
	//}
	//
	//cmd.AddCommand(sendtx.Command())
	//cmd.AddCommand(call.Command())
	//cmd.AddCommand(checktx.Command())

	viper.AutomaticEnv()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func executeCmd(_ *cobra.Command, _ []string) error {

	//if viper.IsSet(versionFlag) {
	//	if version != "" {
	//		fmt.Print(version)
	//	} else {
	//		fmt.Println(branch)
	//	}
	//}
	//
	//if viper.IsSet(buildInfoFlag) {
	//	fmt.Println("commit-hash:", commitHash, " build-time:", buildTime)
	//}

	fmt.Println(version, commitHash, buildTime, branch)

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

	d, _ := json.MarshalIndent(manifest, "", "   ")
	fmt.Println((string)(d))

	toolsStorage := infra.NewDefaultToolsStorage("../../tool-specs")
	toolsStorageAdapter := adapters.NewDefaultStorageAdapter(toolsStorage)
	theChecker := manifestchecker.NewDefaultManifestChecker(toolsStorageAdapter)

	err = theChecker.Check(*manifest, func(msg manifestchecker.Notification) {
		fmt.Printf("%s found:%v version_found:%s version_requeted:%s version_valid:%v", msg.Command, msg.CommandFound, msg.VersionFound, msg.VersionRequested, msg.VersionValid)
	})
	if err != nil {
		panic(err)
	}

	return nil
}
