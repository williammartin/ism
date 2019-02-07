package commands

type RootCommand struct {
	BrokerCommand BrokerCommand `command:"broker" long-description:"The broker command group lets you register, update and deregister Service Brokers from the marketplace"`
}
