# godriver

test and load test driver written in go


## description

godriver is a simple library for building functional and load test oriented test suites
for web applications.


```go
	err := Get("https://www.golang.org").
		HasContentType("text/html").
		SelectorContains("div.rootHeading", "Try Java").
		Exec(nil)

	fmt.Println(err.Error())
	// Output: selection does not contain "Try Java", but was: "Try Go"
```
