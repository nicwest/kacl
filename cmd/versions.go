package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//var Source string

// versionsCmd represents the versions command
var versionsCmd = &cobra.Command{
	Use:     "versions",
	//Aliases: []string{"s"},
	Short:   "List all versions in a changelog",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			fmt.Println("Cannot read changelog")
			return
		}
		fmt.Println(contents.ChangeLogVersions())
	},
}

func init() {
	RootCmd.AddCommand(versionsCmd)
}
