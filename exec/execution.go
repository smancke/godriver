package exec

import (
	"fmt"
	"strings"
	"time"
)

type Execution struct {
	start    time.Time
	end      time.Time
	jobTitle string
	err      error
	context  Context
	retries  uint
}

func StartExecution(jobTitle string, context *Context) *Execution {
	return &Execution{
		start:    time.Now(),
		jobTitle: jobTitle,
		context:  *context,
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
	msgParts := []string{}
	msgParts = append(msgParts, execution.Duration().String(), fmt.Sprintf("%v:", execution.jobTitle))
	if execution.err != nil {
		msgParts = append(msgParts, fmt.Sprintf("%v", execution.err))
	}
	if execution.context.Retries() > 0 {
		msgParts = append(msgParts, fmt.Sprintf("(%v retries)", execution.context.Retries()))
	}
	msgParts = append(msgParts, execution.context.CorrelationId())

	return strings.Join(msgParts[:], " ")
}
