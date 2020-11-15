package cmd

import (
	"github.com/spf13/cobra"
)

// addedCmd represents the added command
var addedCmd = &cobra.Command{
	Use:     "added",
	Aliases: []string{"a"},
	Short:   "Add a change to the list of Unreleased additions",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			return
		}
		addLine(&contents.Unreleased.Added, args)
		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(addedCmd)
}
