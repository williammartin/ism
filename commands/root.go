package commands

type RootCommand struct {
	BrokerCommand   NullCommand     `command:"broker" long-description:"The broker command group lets you register, update and deregister service brokers from the marketplace"`
	ServiceCommand  NullCommand     `command:"service" long-description:"The service command group lets you list the available services in the marketplace."`
	InstanceCommand InstanceCommand `command:"instance" long-description:"The instance command group lets you list, create, update and delete service instances"`
	BindingCommand  BindingCommand  `command:"binding" long-description:"The binding command group lets you list, create, get and delete service bindings"`
}
