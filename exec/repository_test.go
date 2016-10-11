package exec

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
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

func newMockErrorExecution(scenarioName string) Exec {
	return F(scenarioName, func() error {
		return errors.New("Wanted error on \"" + scenarioName + "\"")
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

func Test_ErrorExecutions_Nil(t *testing.T) {
	repo := NewRepository()
	assert.Nil(t, repo.GetErrorExecutions())
}

func Test_ErrorExecutions_Empty(t *testing.T) {
	repo := NewRepository()
	repo.Add(NewTestScenario("", newMockExec("Test_ErrorExecutionsEmpty"), newChannelFactory()), "groupppp", 0)
	repo.RunTestScenarios("groupppp", "")

	assert.Empty(t, repo.GetErrorExecutions())
}

func Test_ErrorExecutions_NotEmpty(t *testing.T) {
	repo := NewRepository()
	repo.Add(NewTestScenario("", newMockErrorExecution("Test_ErrorExecutionsNotEmpty"),
		newChannelFactory()), "ggggroup", 0)
	repo.RunTestScenarios("ggggroup", "")

	assert.NotEmpty(t, repo.GetErrorExecutions())
}
