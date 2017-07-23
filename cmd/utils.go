package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nicwest/kacl/changelog"
)

func getContents() (*changelog.Contents, bool) {
	if _, err := os.Stat("./CHANGELOG.md"); os.IsNotExist(err) {
		color.Red("CHANGELOG.md doesn't exists!")
		fmt.Printf("you can use ")
		color.Set(color.FgCyan)
		fmt.Printf("kacl init")
		color.Unset()
		fmt.Printf("to start a new one\n")
		return nil, false
	}

	f, err := os.Open("./CHANGELOG.md")

	if err != nil {
		color.Red(err.Error())
		return nil, false
	}

	contents, err := changelog.Parse(f)
	f.Close()

	if err != nil {
		color.Red(err.Error())
		return nil, false
	}
	return contents, true
}

func writeContents(contents *changelog.Contents) bool {

	f, err := os.OpenFile("./CHANGELOG.md", os.O_WRONLY, 0644)
	if err != nil {
		color.Red(err.Error())
		return false
	}
	defer f.Close()
	_, err = contents.WriteTo(f)
	if err != nil {
		color.Red(err.Error())
		return false
	}
	return true
}

func addLine(section *string, args []string) {
	line := strings.Join(args, " ")
	if *section == "" {
		*section = fmt.Sprintf("- %s", line)
		return
	}
	*section = fmt.Sprintf("%s\n- %s", *section, line)
}
