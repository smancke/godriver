package exec

import (
	"fmt"
	"time"
)

type Execution struct {
	start         time.Time
	end           time.Time
	jobTitle      string
	err           error
	correlationId string
}

func StartExecution(jobTitle string, correlationId string) *Execution {
	return &Execution{
		start:         time.Now(),
		jobTitle:      jobTitle,
		correlationId: correlationId,
	}
}

func (execution *Execution) End(err error) {
	execution.end = time.Now()
	execution.err = err
}

func (execution *Execution) Duration() time.Duration {
	return execution.end.Sub(execution.start)
}

func (execution *Execution) String() string {
	if execution.err == nil {
		return fmt.Sprintf("%v %v %v", execution.Duration(), execution.jobTitle, execution.correlationId)
	} else {
		return fmt.Sprintf("%v %v: %v %v", execution.Duration(), execution.jobTitle, execution.err, execution.correlationId)
	}
}
