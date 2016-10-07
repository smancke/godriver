package exec

type TestScenario struct {
	Name                  string
	Exec                  Exec
	ContextChannelFactory func() chan Context
}

func NewTestScenario(name string, exec Exec, contextChannelFactory func() chan Context) *TestScenario {
	return &TestScenario{
		Name: name,
		Exec: exec,
		ContextChannelFactory: contextChannelFactory,
	}
}
