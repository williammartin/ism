package commands

type ServiceCommand struct {
	ServiceListCommand   ServiceListCommand `command:"list" `
	CreateServiceCommand NullCommand        `command:"create"`
	UpdateServiceCommand NullCommand        `command:"update"`
	DeleteServiceCommand NullCommand        `command:"delete"`
}

type ServiceListCommand struct {
	UI UI
}

func (cmd *ServiceListCommand) Execute([]string) error {
	headers := []string{"NAME", "PLANS", "BROKER", "DESCRIPTION"}
	data := [][]string{headers}

	data = append(data, []string{"mysql", "small, mid, large", "mysql-broker", "A MySQL Database."})
	data = append(data, []string{"redis", "free, 256mb, 512mb", "redis-broker", "A Redis cache."})
	data = append(data, []string{"overview-service", "simple, complex", "overview-broker", "A test service."})

	cmd.UI.DisplayTable(data)

	return nil
}
