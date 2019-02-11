package commands

import "github.com/pivotal-cf/ism/actors"

//go:generate counterfeiter . ServicesActor

// TODO: godoc
type ServicesActor interface {
	GetServices() ([]actors.Service, error)
}

//go:generate counterfeiter . UI

// TODO: godoc
type UI interface {
	DisplayText(text string, data ...map[string]interface{})
	// DisplayTable(table [][]string)
}

// TODO: godoc
type ServicesCommand struct {
	ListCommand ListCommand `command:"list" long-description:"List the services that are available in the marketplace."`
}

// TODO: godoc
type ListCommand struct {
	UI            UI
	ServicesActor ServicesActor
}

// TODO: godoc
func (cmd *ListCommand) Execute([]string) error {
	//
	// services, err := cmd.Actor.GetServices()
	// if err != nil {
	// 	return err
	// }
	//
	// if len(services) == 0 {
	cmd.UI.DisplayText("No services found.")
	// 	fmt.Println("No brokers found.")
	// 	return nil
	// }
	//
	// servicesTable := [][]string{}
	// for _, service := range services {
	// 	// append service info into table
	// }
	//
	// UI.DisplayTable(servicesTable)

	return nil
}
