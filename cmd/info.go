package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Source string

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"s"},
	Short:   "List information in a change log",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			fmt.Println("Cannot read changelog")
			return
		}
		fmt.Println(contents.ChangeLogInfo(Source))
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&Source, "tag", "t", "Unreleased", "Specify which tag you wants to display")
}
