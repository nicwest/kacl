package cmd

import (
	"html/template"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nicwest/kacl/prompt"
)

type initCmdConfig struct {
	ProjectURL string
	InitialTag string
}

// V1Template is the default keep a change log v1.0.0 template
const V1Template string = `# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added 
- This CHANGELOG file!


[Unreleased]: {{ .ProjectURL }}/compare/{{ .InitialTag }}...HEAD `

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a CHANGELOG.md file",
	Long: `Initialize will create a new CHANGLE.md file in the 
http://keepachangelog.com/en/1.0.0/ format.`,
	Run: func(cmd *cobra.Command, args []string) {

		force, _ := cmd.LocalFlags().GetBool("force")

		if !force {
			if _, err := os.Stat("./CHANGELOG.md"); !os.IsNotExist(err) {
				color.Red("CHANGELOG.md already exists!")
				return
			}
		}

		var cfg initCmdConfig

		prompt.For("Project URL", &cfg.ProjectURL)
		prompt.ForWithDefault("Initial commit", "0.0.1", &cfg.InitialTag)

		t := template.Must(template.New("version").Parse(V1Template))

		f, err := os.Create("./CHANGELOG.md")
		if err != nil {
			color.Red(err.Error())
			return
		}

		err = t.Execute(f, cfg)
		if err != nil {
			color.Red(err.Error())
			return
		}

		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Will force the creation of the CHANGELOG.md file even if it already exists")
}
