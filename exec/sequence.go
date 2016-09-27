package exec

type SequenceExec struct {
	steps []Exec
	name  string
}

func Seq(name string, steps ...Exec) *SequenceExec {
	return &SequenceExec{
		name:  name,
		steps: steps,
	}
}

func (s *SequenceExec) String() string {
	return s.name
}

func (s *SequenceExec) Exec(cntx Context) error {
	for _, step := range s.steps {
		err := step.Exec(cntx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add a Step to the SequenceExec
func (s *SequenceExec) Add(r Exec) *SequenceExec {
	s.steps = append(s.steps, r)
	return s
}

// Add a Step n times to the SequenceExec
func (s *SequenceExec) AddN(n int, r Exec) *SequenceExec {
	for i := 0; i < n; i++ {
		s.steps = append(s.steps, r)
	}
	return s
}
