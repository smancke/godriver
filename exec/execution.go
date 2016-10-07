package exec

import (
	"fmt"
	"time"
)

type Execution struct {
	start    time.Time
	end      time.Time
	jobTitle string
	err      error
}

func StartExecution(jobTitle string) *Execution {
	return &Execution{
		start:    time.Now(),
		jobTitle: jobTitle,
	}
}

func (execution *Execution) End(err error) {
	execution.end = time.Now()
	execution.err = err
}

func (execution *Execution) Duration() time.Duration {
	return execution.end.Sub(execution.start)
}

func (execution *Execution) Error() error {
	return execution.err
}

func (execution *Execution) String() string {
	if execution.err == nil {
		return fmt.Sprintf("%v %v", execution.Duration(), execution.jobTitle)
	} else {
		return fmt.Sprintf("%v %v: %v", execution.Duration(), execution.jobTitle, execution.err)
	}
}
