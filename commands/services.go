package commands

import (
	"strings"

	"github.com/pivotal-cf/ism/usecases"
)

//go:generate counterfeiter . ListServicesUsecase

// TODO: godoc
type ListServicesUsecase interface {
	GetServices() ([]*usecases.Service, error)
}

//go:generate counterfeiter . UI

// TODO: godoc
type UI interface {
	DisplayText(text string, data ...map[string]interface{})
	DisplayTable(table [][]string)
}

// TODO: godoc
// TODO: Rename to Service
type ServicesCommand struct {
	ListCommand ListCommand `command:"list" long-description:"List the services that are available in the marketplace."`
}

// TODO: godoc
type ListCommand struct {
	UI                  UI
	ListServicesUsecase ListServicesUsecase
}

// TODO: godoc
// TODO: Rename to ServiceList
func (cmd *ListCommand) Execute([]string) error {
	services, err := cmd.ListServicesUsecase.GetServices()
	if err != nil {
		return err
	}

	if len(services) == 0 {
		cmd.UI.DisplayText("No services found.")
		return nil
	}

	data := buildTableData(services)
	cmd.UI.DisplayTable(data)

	return nil
}

func buildTableData(services []*usecases.Service) [][]string {
	headers := []string{"SERVICE", "PLANS", "BROKER", "DESCRIPTION"}
	data := [][]string{headers}

	for _, s := range services {
		row := []string{s.Name, strings.Join(s.PlanNames, ", "), s.BrokerName, s.Description}
		data = append(data, row)
	}

	return data
}
