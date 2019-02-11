package commands

import "fmt"

type ServicesCommand struct {
	ListCommand ListCommand `command:"list" long-description:"List the services that are available in the marketplace."`
}

type ListCommand struct{}

func (cmd *ListCommand) Execute([]string) error {
	fmt.Println("No brokers found.")
	return nil
}
