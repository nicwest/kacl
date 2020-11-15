package cmd

import (
	"github.com/spf13/cobra"
)

// changedCmd represents the changed command
var changedCmd = &cobra.Command{
	Use:     "changed",
	Aliases: []string{"c"},
	Short:   "Add a change to the list of Unreleased changes",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			return
		}
		addLine(&contents.Unreleased.Changed, args)
		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(changedCmd)
}
