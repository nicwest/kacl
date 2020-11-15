package cmd

import (
	"github.com/spf13/cobra"
)

// deprecatedCmd represents the deprecated command
var deprecatedCmd = &cobra.Command{
	Use:     "deprecated",
	Aliases: []string{"d"},
	Short:   "Add a change to the list of Unreleased deprecations",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			return
		}
		addLine(&contents.Unreleased.Deprecated, args)
		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(deprecatedCmd)
}
