package commands

import (
	"github.com/pivotal-cf/ism/osbapi"
)

//go:generate counterfeiter . BrokerRegistrar

//TODO: godoc
type BrokerRegistrar interface {
	Register(*osbapi.Broker) error
}

type BrokerCommand struct {
	RegisterCommand NullCommand `command:"register" long-description:"Register a service broker into the marketplace"`
}

type RegisterCommand struct {
	Name     string `long:"name" description:"Name of the service broker" required:"true"`
	URL      string `long:"url" description:"URL of the service broker" required:"true"`
	Username string `long:"username" description:"Username of the service broker" required:"true"`
	Password string `long:"password" description:"Password of the service broker" required:"true"`

	UI              UI
	BrokerRegistrar BrokerRegistrar
}

func (cmd *RegisterCommand) Execute([]string) error {
	newBroker := &osbapi.Broker{
		Name:     cmd.Name,
		URL:      cmd.URL,
		Username: cmd.Username,
		Password: cmd.Password,
	}

	if err := cmd.BrokerRegistrar.Register(newBroker); err != nil {
		return err
	}

	cmd.UI.DisplayText("Broker '{{.BrokerName}}' registered.", map[string]interface{}{"BrokerName": cmd.Name})

	return nil
}
