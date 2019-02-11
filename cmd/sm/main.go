package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/ui"
)

func main() {
	UI := &ui.UI{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	servicesActor := &actors.ServicesActor{}

	rootCommand := commands.RootCommand{
		BrokerCommand: commands.BrokerCommand{},
		ServicesCommand: commands.ServicesCommand{
			ListCommand: commands.ListCommand{
				UI:            UI,
				ServicesActor: servicesActor,
			},
		},
	}
	parser := flags.NewParser(&rootCommand, flags.HelpFlag|flags.PassDoubleDash)

	_, err := parser.Parse()

	if err != nil {
		_, err := parser.ParseArgs([]string{"--help"})
		fmt.Println(err)
	}
}
