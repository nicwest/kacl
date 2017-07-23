package cmd

import (
	"github.com/spf13/cobra"
)

// removedCmd represents the removed command
var removedCmd = &cobra.Command{
	Use:     "removed",
	Aliases: []string{"r"},
	Short:   "Add a change to the list of Unreleased removals",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			return
		}
		addLine(&contents.Unreleased.Removed, args)
		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(removedCmd)
}
