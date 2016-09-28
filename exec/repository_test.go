package exec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockResult = ""

func mockTestFactory(resultString string) TestFactory {
	return func(cntx Context) (Exec, chan Context) {
		c := make(chan Context, 1)
		c <- cntx
		close(c)
		return F(resultString, func() error {
			mockResult += resultString
			return nil
		}), c
	}
}

func Test_Repository(t *testing.T) {
	a := assert.New(t)

	repo := NewRepository()
	repo.AddTest("group1", "spec11", mockTestFactory("spec11,"), NewDefaultContext(), "foo", "bar")
	repo.AddTest("group1", "spec12", mockTestFactory("spec12,"), NewDefaultContext(), "foo", "bazz")
	repo.AddTest("group2", "spec21", mockTestFactory("spec21,"), NewDefaultContext())

	mockResult = ""
	repo.RunTests("group1", "")
	a.Equal("spec11,spec12,", mockResult)

	mockResult = ""
	repo.RunTests("", "")
	a.Equal("spec11,spec12,spec21,", mockResult)

	mockResult = ""
	repo.RunTests("", "spec12")
	a.Equal("spec12,", mockResult)

	mockResult = ""
	repo.RunTests("", "spec.1")
	a.Equal("spec11,spec21,", mockResult)

	mockResult = ""
	repo.RunTests("", "Spec21")
	a.Equal("", mockResult)

	mockResult = ""
	repo.RunTests("", "", "foo")
	a.Equal("spec11,spec12,", mockResult)

	mockResult = ""
	repo.RunTests("", "", "bazz")
	a.Equal("spec12,", mockResult)
}
