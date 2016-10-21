package exec

type FuncExec struct {
	f    func() error
	name string
}

func F(name string, f func() error) *FuncExec {
	return &FuncExec{
		name: name,
		f:    f,
	}
}

func (s *FuncExec) String(cntx Context) string {
	return cntx.ExpandVarsNoError(s.name)
}

func (s *FuncExec) Exec(cntx Context) error {
	return s.f()
}
