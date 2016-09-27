package exec

import (
	"fmt"
)

func ExampleGet() {
	err := Get("https://www.golang.org").
		HasContentType("text/html").
		SelectorContains("div.rootHeading", "Try Java").
		Exec(NewDefaultContext())

	fmt.Println(err.Error())
	// Output: selection does not contain "Try Java", but was: "Try Go"
}
