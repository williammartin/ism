package commands

type RootCommand struct {
	BrokerCommand   BrokerCommand   `command:"broker" long-description:"The broker command group lets you register, update and deregister Service Brokers from the marketplace"`
	ServicesCommand ServicesCommand `command:"services" long-description:"The services command group lets you list the available services in the marketplace."`
}
