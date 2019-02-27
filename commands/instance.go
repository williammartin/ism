package commands

type InstanceCommand struct {
	InstanceListCommand   InstanceListCommand `command:"list" `
	CreateInstanceCommand NullCommand         `command:"create"`
	UpdateInstanceCommand NullCommand         `command:"update"`
	DeleteInstanceCommand NullCommand         `command:"delete"`
}

type InstanceListCommand struct {
	UI UI
}

func (cmd *InstanceListCommand) Execute([]string) error {
	headers := []string{"INSTANCE NAME", "SERVICE", "PLAN", "CREATED AT"}
	data := [][]string{headers}

	data = append(data, []string{"my-db", "mysql", "small", "2019-01-22-17:43"})
	data = append(data, []string{"redis-tiny", "redis", "free", "2017-10-07-03:58"})
	data = append(data, []string{"identity-cache", "redis", "256mb", "2017-01-10-06:21"})
	data = append(data, []string{"identity-db", "mysql", "large", "2018-04-16-15:48"})
	data = append(data, []string{"DB-FOR-PAYMENTS", "mysql", "large", "2018-04-02-14:42"})
	data = append(data, []string{"my-db2-12345", "mysql", "mid", "2017-11-28-04:53"})
	data = append(data, []string{"6fc2f003e130-31", "redis", "free", "2018-10-06-02:35"})
	data = append(data, []string{"identity-db2", "mysql", "large", "2019-02-20-04:22"})
	data = append(data, []string{"TEST-DB", "mysql", "large", "2018-06-18-12:45"})

	cmd.UI.DisplayTable(data)

	return nil
}
