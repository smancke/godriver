package exec

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Context interface {
	// TestNumber is the numer of the actual test.
	TestNumber() int

	// ExecutionCount is the total number of test executions for a given test.
	ExecutionCount() int

	// ExecutionConcurrency is the number of parallel workers, which should be used.
	ExecutionConcurrency() int

	// Test returns the configuration data
	// for a single test interation.
	Test() map[string]string

	// Env returns the environment configuration which is common for all
	// test groups.
	Env() map[string]string

	// ExpandVars executes the supplied go template with the context as data context
	ExpandVars(template string) (string, error)

	// Derive creates a copy of the context, where the test data is
	// field wise overwritten by the supplied test data and the
	// test number is incremented.
	Derive(overrideValues map[string]string) Context

	// Populate can be used to create test data for the number of ExecutionCount tests.
	// It calls the supplied closure for each test and derives a new Context using the test data returned by the supplied function.
	// The creation is done in a go routine and supplied over returned channel.
	// The channel will be closed afer sending the last entry.
	Populate(createTestDataClosure func(testNumber int) map[string]string) chan Context

	// Returns the CorrelationId
	// CorrelationId became a http header attribute
	CorrelationId() string
}

type ContextImpl struct {
	test                 map[string]string
	env                  map[string]string
	testNumber           int
	executionCount       int
	executionConcurrency int
	correlationId        string
}

// NewDefaultContext creates a new context without data and
// - executionCount = 1
// - executionConcurrency = 1
func NewDefaultContext() *ContextImpl {
	return &ContextImpl{
		env:                  make(map[string]string),
		test:                 make(map[string]string),
		testNumber:           0,
		executionCount:       1,
		executionConcurrency: 1,
		correlationId:        "",
	}
}

// NewContext creates a new context
// - executionCount number of executions
// - executionConcurrency number of parallel workers
// - env base data, which may be nil
func NewContext(executionCount, executionConcurrency int, env map[string]string) *ContextImpl {
	cntx := &ContextImpl{
		env:                  env,
		test:                 make(map[string]string),
		testNumber:           0,
		executionCount:       executionCount,
		executionConcurrency: executionConcurrency,
		correlationId:        "",
	}
	if cntx.env == nil {
		cntx.env = make(map[string]string)
	}
	return cntx
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (cntx *ContextImpl) ExecutionCount() int {
	return cntx.executionCount
}

func (cntx *ContextImpl) ExecutionConcurrency() int {
	return cntx.executionConcurrency
}

func (cntx *ContextImpl) Env() map[string]string {
	return cntx.env
}

func (cntx *ContextImpl) Test() map[string]string {
	return cntx.test
}

func (cntx *ContextImpl) TestNumber() int {
	return cntx.testNumber
}

func (cntx *ContextImpl) CorrelationId() string {
	return cntx.correlationId
}

func (cntx *ContextImpl) ExpandVars(tpl string) (string, error) {
	t, err := template.New("template").Parse(tpl)
	if err != nil {
		return "", err
	}
	b := bytes.NewBuffer(nil)
	err = t.ExecuteTemplate(b, "template", cntx)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (cntx *ContextImpl) Derive(overrideValues map[string]string) Context {
	contextCopy := *cntx
	contextCopy.testNumber++
	contextCopy.correlationId = randStringBytes(10)
	contextCopy.test = make(map[string]string)
	for k, v := range cntx.test {
		contextCopy.test[k] = v
	}
	for k, v := range overrideValues {
		contextCopy.test[k] = v
	}
	return &contextCopy
}

func (cntx *ContextImpl) Populate(createTestDataClosure func(testNumber int) map[string]string) chan Context {
	resultChannel := make(chan Context)
	go func() {
		var currentContext Context
		currentContext = cntx
		for i := cntx.TestNumber() + 1; i <= cntx.ExecutionCount(); i++ {
			currentContext = currentContext.Derive(createTestDataClosure(i))
			resultChannel <- currentContext
		}
		close(resultChannel)
	}()
	return resultChannel
}
