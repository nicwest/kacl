package cmd

import (
	"bytes"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nicwest/kacl/changelog"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release [tag]",
	Short: "Create a new release",
	Long: `Create a new release by moving the current Unreleased changes into 
a new change with the given tag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			color.Red("no tag specified!")
			return
		}

		tag := args[0]

		contents, ok := getContents()
		if !ok {
			return
		}

		changes := contents.Unreleased
		changes.Tag = tag
		changes.Time = time.Now()

		rest := bytes.NewBufferString("")
		changes.WriteTo(rest)
		rest.WriteString("\n")
		rest.WriteString(contents.Rest)
		contents.Rest = rest.String()

		if len(contents.Refs) > 0 {

			lastTag := "TAIL"
			if len(contents.Changes) > 1 {
				last := contents.Changes[1]
				lastTag = last.Tag
			}

			base := contents.Refs[0].BaseURL
			contents.Refs = append([]changelog.Reference{
				{
					Tag:     "Unreleased",
					From:    tag,
					To:      "HEAD",
					BaseURL: base,
				},
				{
					Tag:     tag,
					From:    lastTag,
					To:      tag,
					BaseURL: base,
				},
			}, contents.Refs[1:]...)
		}

		contents.Unreleased = changelog.NewChanges("Unreleased")

		writeContents(contents)
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
