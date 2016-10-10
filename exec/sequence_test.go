package exec

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Sequence(t *testing.T) {
	a := assert.New(t)
	cntx := NewDefaultContext()

	result := make(chan string, 3)

	err := Seq("some test",
		F("Test a", func() error {
			result <- "a"
			return nil
		}),
		F("Test b", func() error {
			result <- "b"
			return nil

		}),
		F("Test c", func() error {
			result <- "c"
			return errors.New("c has an error")
		}),
		F("Test d", func() error {
			result <- "d"
			return nil
		})).Exec(cntx)

	// never executed d
	a.Equal(3, len(result))
	a.Equal("a", <-result)
	a.Equal("b", <-result)
	a.Equal("c", <-result)
	a.Error(err)
	a.Equal("c has an error", err.Error())
}
