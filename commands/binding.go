package commands

type BindingCommand struct {
	BindingListCommand   BindingListCommand `command:"list" long-description:"List service bindings that have been created.\n\nFor Cloud Foundry, the LOCATION column will show platform-name/org-name/space-name/app-name."`
	CreateBindingCommand NullCommand        `command:"create"`
	GetBindingCommand    NullCommand        `command:"get"`
	DeleteBindingCommand NullCommand        `command:"delete"`
}

type BindingListCommand struct {
	UI UI
}

func (cmd *BindingListCommand) Execute([]string) error {
	headers := []string{"BINDING NAME", "INSTANCE NAME", "CREATED AT", "LOCATION"}
	data := [][]string{
		headers,
		[]string{"id-mservice-to-db-g", "identity-cache", "2018-03-18-21:50", "cf-west-01/dev-green/lob-ops/identity"},
		[]string{"depr-test-bind", "my-db", "2018-05-30-13:03", "cf-east-02/sandbox/temp/test-app-dora"},
		[]string{"a94a8fe5ccb19b", "identity-db2", "2017-06-07-15:22", "cf-01-01/accounting/test/1e987982fbbd3-test"},
		[]string{"binding-0187t5bq", "6fc2f003e130-31", "2016-08-12-05:05", "cf-east-02/sandbox/temp/55db843d939"},
		[]string{"identity-cache", "identity-cache", "2017-08-24-15:13", "cf-west-01/dev-green/live/identity"},
		[]string{"TEST-APP-BIND", "TEST-DB", "2016-05-03-07:50", "cf-west-01/dev-green/lob-ops/TEST"},
		[]string{"mysql-bind-2", "my-db2-12345", "2018-05-14-07:52", "cf-01-01/user-id/space-team/identity"},
		[]string{"db-binding", "DB-FOR-PAYMENTS", "2018-04-19-14:12", "cf-01-03/accounting/prod/1e987982fbbd3-prod"},
		[]string{"id-mservice-to-db-b", "identity-cache", "2018-03-04-06:12", "cf-west-01/dev-blue/lob-ops/identity"},
	}

	cmd.UI.DisplayTable(data)

	return nil
}
