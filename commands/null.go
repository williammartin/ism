package commands

import "errors"

type UI interface {
	DisplayText(text string, data ...map[string]interface{})
	DisplayTable(table [][]string)
}

type NullCommand struct{}

func (cmd *NullCommand) Execute([]string) error {
	return errors.New("NOT INCLUDED IN PROTOTYPE")
}
