package commands

type BrokerCommand struct {
	BrokerListCommand   BrokerListCommand `command:"list" `
	CreateBrokerCommand NullCommand       `command:"create"`
	UpdateBrokerCommand NullCommand       `command:"update"`
	DeleteBrokerCommand NullCommand       `command:"delete"`
}

type BrokerListCommand struct {
	UI UI
}

func (cmd *BrokerListCommand) Execute([]string) error {
	headers := []string{"NAME", "URL"}
	data := [][]string{headers}

	data = append(data, []string{"mysql-broker", "http://10.2.0.10"})
	data = append(data, []string{"redis-broker", "http://10.0.0.44"})
	data = append(data, []string{"overview-broker", "http://10.3.33.8"})

	cmd.UI.DisplayTable(data)

	return nil
}
