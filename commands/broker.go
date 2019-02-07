package commands

import "fmt"

type BrokerCommand struct {
}

type RegisterCommand struct {
	Name     string `long:"name" description:"name of the broker to regsiter"`
	URL      string `long:"url" description:"url of the broker to regsiter"`
	Username string `long:"username" description:"username of the broker to regsiter"`
	Password string `long:"password" description:"password of the broker to regsiter"`
}

func (cmd *RegisterCommand) Execute([]string) error {
	fmt.Printf("Broker '%s' registered.\n", cmd.Name)
	return nil
}
