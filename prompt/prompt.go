package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func scanLine() (string, error) {
	input := bufio.NewReader(os.Stdin)
	line, err := input.ReadString('\n')
	if err != nil {
		return line, err
	}

	return strings.TrimSuffix(line, "\n"), nil
}

func For(p string, v *string) {
	color.Set(color.FgCyan)
	fmt.Printf("%s: ", p)
	color.Unset()

	in, err := scanLine()

	if err != nil {
		color.Red(err.Error())
		return
	}

	*v = in
}

func ForWithDefault(p string, d string, v *string) {
	color.Set(color.FgCyan)
	fmt.Printf("%s ", p)
	color.Set(color.FgMagenta)
	fmt.Printf("[%s]: ", d)
	color.Unset()

	in, err := scanLine()

	if err != nil {
		color.Red(err.Error())
		return
	}

	if in != "" {
		*v = in
		return
	}

	*v = d
}
