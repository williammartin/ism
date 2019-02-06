package main

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

type smCommand struct {
	BrokerCommand brokerCommand `command:"broker"`
}

type brokerCommand struct {
	RegisterCommand registerCommand `command:"register"`
}

type registerCommand struct {
	Name string `long:"name" description:"name of the broker to regsiter"`
}

func (cmd *registerCommand) Execute([]string) error {
	fmt.Println("listyList")
	return nil
}

func main() {
	cmd := &smCommand{}
	parser := flags.NewParser(cmd, flags.HelpFlag|flags.PassDoubleDash)

	fmt.Println("CLI to interact with the Services Marketplace")
	_, err := parser.Parse()

	if err != nil {
		_, err := parser.ParseArgs([]string{"--help"})
		fmt.Println(err)
	}
}
