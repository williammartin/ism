package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/ui"
)

func main() {
	UI := &ui.UI{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	rootCommand := commands.RootCommand{
		InstanceCommand: commands.InstanceCommand{
			InstanceListCommand: commands.InstanceListCommand{
				UI: UI,
			},
		},
		BindingCommand: commands.BindingCommand{
			BindingListCommand: commands.BindingListCommand{
				UI: UI,
			},
		},
		BrokerCommand: commands.BrokerCommand{
			BrokerListCommand: commands.BrokerListCommand{
				UI: UI,
			},
		},
		ServiceCommand: commands.ServiceCommand{
			ServiceListCommand: commands.ServiceListCommand{
				UI: UI,
			},
		},
	}
	parser := flags.NewParser(&rootCommand, flags.HelpFlag|flags.PassDoubleDash)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	if os.Args[1] == "help" {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	_, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
	}
}
