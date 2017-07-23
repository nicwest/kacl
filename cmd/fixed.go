package cmd

import (
	"github.com/spf13/cobra"
)

// fixedCmd represents the fixed command
var fixedCmd = &cobra.Command{
	Use:     "fixed",
	Aliases: []string{"f"},
	Short:   "Add a change to the list of Unreleased fixes",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			return
		}
		addLine(&contents.Unreleased.Fixed, args)
		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(fixedCmd)
}
