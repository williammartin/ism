package main

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/ism/commands"
)

func main() {
	parser := flags.NewParser(&commands.RootCommand{}, flags.HelpFlag|flags.PassDoubleDash)

	_, err := parser.Parse()

	if err != nil {
		_, err := parser.ParseArgs([]string{"--help"})
		fmt.Println(err)
	}
}
