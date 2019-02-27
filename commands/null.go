package commands

import "errors"

type NullCommand struct{}

func (cmd *NullCommand) Execute([]string) error {
	return errors.New("NOT INCLUDED IN PROTOTYPE")
}
