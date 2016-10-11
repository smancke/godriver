package exec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockResult = ""

func newChannelFactory() func() chan Context {
	return func() chan Context {
		c := make(chan Context, 1)
		c <- NewDefaultContext()
		close(c)
		return c
	}
}

func newMockExec(scenarioName string) Exec {
	return F(scenarioName, func() error {
		if len(mockResult) == 0 {
			mockResult = scenarioName
		} else {
			mockResult += "," + scenarioName
		}
		return nil
	})
}

func Test_Repository(t *testing.T) {
	a := assert.New(t)

	repo := NewRepository()
	repo.Add(NewTestScenario("spec11", newMockExec("spec11"), newChannelFactory()),
		"group1", 1, "foo", "bar")

	repo.Add(NewTestScenario("spec12", newMockExec("spec12"), newChannelFactory()),
		"group1", 1, "foo", "bazz")

	repo.Add(NewTestScenario("spec21", newMockExec("spec21"), newChannelFactory()),
		"group2", 1)

	mockResult = ""
	repo.RunTestScenarios("group1", "")
	a.Equal("spec11,spec12", mockResult)

	mockResult = ""
	repo.RunTestScenarios("", "")
	a.Equal("spec11,spec12,spec21", mockResult)

	mockResult = ""
	repo.RunTestScenarios("", "spec12")
	a.Equal("spec12", mockResult)

	mockResult = ""
	repo.RunTestScenarios("", "spec.1")
	a.Equal("spec11,spec21", mockResult)

	mockResult = ""
	repo.RunTestScenarios("", "Spec21")
	a.Equal("", mockResult)

	mockResult = ""
	repo.RunTestScenarios("", "", "foo")
	a.Equal("spec11,spec12", mockResult)

	mockResult = ""
	repo.RunTestScenarios("", "", "bazz")
	a.Equal("spec12", mockResult)
}
