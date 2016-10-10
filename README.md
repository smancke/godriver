# godriver

[![Build Status](https://api.travis-ci.org/smancke/godriver.svg)](https://travis-ci.org/smancke/godriver) [![Coverage Status](https://coveralls.io/repos/smancke/godriver/badge.svg?branch=master&service=github)](https://coveralls.io/github/smancke/godriver?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/smancke/godriver)](https://goreportcard.com/report/github.com/smancke/godriver)

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
