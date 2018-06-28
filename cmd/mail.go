// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/nicwest/kacl/prompt"
	"github.com/spf13/cobra"
)

type ReleaseInfo struct {
	TimeString string
	DateString string
	Version    string
	Build      string
	Notes      string
	Component  string
	Name       string
}

// mailCmd represents the mail command
var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Create a email from changelog",
	Run: func(cmd *cobra.Command, args []string) {
		contents, ok := getContents()
		if !ok {
			fmt.Println("Cannot read changelog")
			return
		}
		fmt.Println(generateEmail(contents.ChangeLogInfo("Unreleased")))
	},
}

func generateEmail(changeLogInfo string) string {

	var BuildNumber string
	var Name string
	var ReleaseDate string
	var VersionNumber string
	var ReleaseComponent string

	prompt.For("Name: ", &Name)
	prompt.For("Release Component: ", &ReleaseComponent)
	prompt.For("Release Date (yyyy/MM/dd HH:mm): ", &ReleaseDate)
	prompt.For("Version Number: ", &VersionNumber)
	prompt.For("Build Number: ", &BuildNumber)

	emailTmpl := `
Hi All,
 
This is to inform you that starting at {{.TimeString}} today, {{.DateString}}, we’ll be making a deployment of version {{.Version}} (build #{{.Build}}) of {{.Component}}.
 
Release Notes:
{{.Notes}}
 
Deployment Owner: {{.Name}}
Components that will be released: {{.Component}}
 
Best Regards,
{{.Name}}
	`
	ansic := "2006/01/02 15:04"
	t1, err := time.Parse(ansic, ReleaseDate)

	if err != nil {
		panic(err)
	}

	timeString := t1.Add(time.Hour*time.Duration(2)).Format("15:04") + "(HKT " + t1.Format("15:04") + ")"
	info := ReleaseInfo{TimeString: timeString, DateString: t1.Format("Jan 02 2006"), Version: VersionNumber, Build: BuildNumber, Notes: changeLogInfo, Name: Name, Component: ReleaseComponent}

	tmpl, err := template.New("test").Parse(emailTmpl)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, info)
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func init() {
	RootCmd.AddCommand(mailCmd)

}
