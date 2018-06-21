package cmd

import (
	"fmt"
	"strings"

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
			return
		}
		for _, element := range contents.Changes {
			if strings.Compare(element.Tag, Source) == 0 {
				fmt.Println(element.Tag)
				fmt.Println()

				if element.Added != "" {
					fmt.Println("Added")
					fmt.Println(element.Added)
					fmt.Println()
				}

				if element.Changed != "" {
					fmt.Println("Changed")
					fmt.Println(element.Changed)
					fmt.Println()
				}

				if element.Deprecated != "" {
					fmt.Println("Deprecated")
					fmt.Println(element.Deprecated)
					fmt.Println()
				}

				if element.Fixed != "" {
					fmt.Println("Fixed")
					fmt.Println(element.Fixed)
					fmt.Println()
				}

				if element.Removed != "" {
					fmt.Println("Removed")
					fmt.Println(element.Removed)
					fmt.Println()
				}

				if element.Security != "" {
					fmt.Println("Security")
					fmt.Println(element.Security)
					fmt.Println()
				}

				return
			}
		}
		fmt.Println("Version Not Found")
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&Source, "tag", "t", "Unreleased", "Specify which tag you wants to display")
}
