package commands

import "fmt"

type BrokerCommand struct {
	RegisterCommand RegisterCommand `command:"register" long-description:"Register a Service Broker into the marketplace"`
}

type RegisterCommand struct {
	Name     string `long:"name" description:"Name of the Service Broker" required:"true"`
	URL      string `long:"url" description:"URL of the Service Broker" required:"true"`
	Username string `long:"username" description:"Username of the Service Broker"`
	Password string `long:"password" description:"Password of the Service Broker"`
}

func (cmd *RegisterCommand) Execute([]string) error {
	fmt.Printf("Broker '%s' registered.\n", cmd.Name)
	return nil
}
