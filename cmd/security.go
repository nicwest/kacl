package cmd

import (
	"github.com/spf13/cobra"
)

// securityCmd represents the security command
var securityCmd = &cobra.Command{
	Use:     "security",
	Aliases: []string{"s"},
	Short:   "Add a change to the list of Unreleased security updates",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			return
		}
		addLine(&contents.Unreleased.Security, args)
		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(securityCmd)
}
