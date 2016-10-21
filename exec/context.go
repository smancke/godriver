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
	// TestNumber is the number of the actual test.
	TestNumber() int

	// Test returns the configuration data
	// for a single test integration.
	Test() map[string]string

	// Env returns the environment configuration which is common for all
	// test groups.
	Env() map[string]string

	// ExpandVars executes the supplied go template with the context as data context
	ExpandVars(template string) (string, error)

	// ExpandVarsNoError Same as expand vars, but returning template it self in case of an error
	ExpandVarsNoError(template string) string

	// Derive creates a copy of the context, where the test data is
	// field wise overwritten by the supplied test data and the
	// test number is incremented.
	Derive(overrideValues map[string]string) Context

	// Populate can be used to create test data for the number of ExecutionCount tests.
	// It calls the supplied closure for each test and derives a new Context using the test data returned by the supplied function.
	// The creation is done in a go routine and supplied over returned channel.
	// The channel will be closed after sending the last entry.
	Populate(n int, createTestDataClosure func(testNumber int) map[string]string) chan Context

	// CorrelationId is the id which should be transferred in the service chain
	CorrelationId() string
	SetRetries(retries uint)
	Retries() uint
}

type ContextImpl struct {
	test          map[string]string
	env           map[string]string
	testNumber    int
	correlationId string
	retries       uint
}

// NewDefaultContext creates a new context without data
func NewDefaultContext() *ContextImpl {
	return &ContextImpl{
		env:           make(map[string]string),
		test:          make(map[string]string),
		testNumber:    0,
		correlationId: "",
		retries:       0,
	}
}

// NewContext creates a new context
// - env base data, which ma be nil
func NewContext(env map[string]string) *ContextImpl {
	cntx := &ContextImpl{
		env:           env,
		test:          make(map[string]string),
		testNumber:    0,
		correlationId: "",
		retries:       0,
	}
	if cntx.env == nil {
		cntx.env = make(map[string]string)
	}
	return cntx
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

func randStringBytes(n int) string {
	if n < 1 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (cntx *ContextImpl) ExpandVarsNoError(tpl string) string {
	result, err := cntx.ExpandVars(tpl)
	if err != nil {
		return tpl
	}
	return result
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

func (cntx *ContextImpl) Populate(n int, createTestDataClosure func(testNumber int) map[string]string) chan Context {
	resultChannel := make(chan Context)
	go func() {
		var currentContext Context
		currentContext = cntx
		for i := cntx.TestNumber() + 1; i <= n; i++ {
			currentContext = currentContext.Derive(createTestDataClosure(i))
			resultChannel <- currentContext
		}
		close(resultChannel)
	}()
	return resultChannel
}

func (cntx *ContextImpl) Retries() uint {
	return cntx.retries
}

func (cntx *ContextImpl) SetRetries(retries uint) {
	cntx.retries = retries
}
