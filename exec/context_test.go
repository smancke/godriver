package exec

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_Context_Derive(t *testing.T) {
	a := assert.New(t)

	cntx := NewDefaultContext()
	cntx.Test()["a"] = "b"
	cntx.Test()["foo"] = "bar"

	a.Equal(0, cntx.TestNumber())

	derived := cntx.Derive(map[string]string{
		"foo": "bazz",
	})

	a.NotEmpty(derived.CorrelationId())
	a.Equal(0, cntx.TestNumber())
	a.Equal(1, derived.TestNumber())
	a.Equal("b", derived.Test()["a"])
	a.Equal("bazz", derived.Test()["foo"])
}

func Test_Context_ExpandVars(t *testing.T) {
	a := assert.New(t)

	cntx := NewContext(map[string]string{
		"host": "example.com",
	})
	cntx.Test()["id"] = "4711"

	result, err := cntx.ExpandVars("http://{{.Env.host}}/id={{.Test.id}}")

	a.NoError(err)
	a.Equal("http://example.com/id=4711", result)

	_, err = cntx.ExpandVars("{{.Foo}}")
	a.Error(err)

	_, err = cntx.ExpandVars("{{}")
	a.Error(err)
}

func Test_Context_Populate(t *testing.T) {
	a := assert.New(t)

	cntx := NewContext(nil)
	testContextChannel := cntx.Populate(3, func(testNumber int) map[string]string {
		return map[string]string{
			"contextNumber": strconv.Itoa(testNumber),
		}
	})

	cntx1 := <-testContextChannel
	a.Equal(1, cntx1.TestNumber())
	a.Equal("1", cntx1.Test()["contextNumber"])

	cntx2 := <-testContextChannel
	a.Equal(2, cntx2.TestNumber())
	a.Equal("2", cntx2.Test()["contextNumber"])

	cntx3 := <-testContextChannel
	a.Equal(3, cntx3.TestNumber())
	a.Equal("3", cntx3.Test()["contextNumber"])

	_, moreItems := <-testContextChannel
	a.False(moreItems)
}

func Test_RandStringBytes(t *testing.T) {
	a := assert.New(t)

	correlationId1 := randStringBytes(10)
	correlationId2 := randStringBytes(10)
	correlationId3 := randStringBytes(-1)
	correlationId4 := randStringBytes(0)

	a.NotEqual(correlationId1, correlationId2)

	a.Equal(10, len(correlationId1))
	a.Equal(10, len(correlationId2))
	a.Equal(0, len(correlationId3))
	a.Equal(0, len(correlationId4))
}
