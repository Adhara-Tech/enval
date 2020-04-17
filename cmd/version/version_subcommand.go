package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "SNAPSHOT"
var branch = ""
var commitHash = ""
var buildTime = ""

func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:  "version",
		Long: "prints version",
		RunE: executeVersionCmd,
	}

	return cmd
}

func executeVersionCmd(_ *cobra.Command, _ []string) error {

	if version != "" {
		fmt.Print(version)
	} else {
		fmt.Print(branch)
	}

	fmt.Println(" commit-hash:", commitHash, " build-time:", buildTime)

	return nil
}
