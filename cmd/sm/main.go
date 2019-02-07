package main

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
	"github.com/pivotal-cf/ism/commands"
)

func main() {
	parser := flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)

	// TODO: Do we need empty type for BrokerCommand?
	brokerCommand, _ := parser.AddCommand("broker", "broker commands", "The broker command group lets you register, update and deregister service brokers from the marketplace", &commands.BrokerCommand{})

	brokerCommand.AddCommand("register", "register a broker", "Register a service broker into the marketplace", &commands.RegisterCommand{})

	_, err := parser.Parse()

	if err != nil {
		_, err := parser.ParseArgs([]string{"--help"})
		fmt.Println(err)
	}
}
