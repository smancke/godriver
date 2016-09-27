package exec

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var html = `<html>
<head>
</head>
<body>
  <h1>hello world</h1>
  <div id="foo">bar</div>
  <div class="item">title</div>
</body>
<html>`

func Test_Http_Get(t *testing.T) {
	a := assert.New(t)
	cntx := ExecutionContextImpl{}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(html))
	}))
	defer server.Close()

	a.Error(Get("h :// invalid").Exec(cntx))
	a.NoError(Get(server.URL).Exec(cntx))

	// Content-Tyoe
	a.NoError(Get(server.URL).HasContentType("text/html").Exec(cntx))
	a.Error(Get(server.URL).HasContentType("text/plain").Exec(cntx))

	// Contains
	a.NoError(Get(server.URL).Contains("hello world").Exec(cntx))

	// goquery
	a.NoError(Get(server.URL).
		SelectorContains("body", "world").
		Exec(cntx))
	a.NoError(Get(server.URL).
		SelectorContains("#foo", "ba").
		Exec(cntx))
	a.NoError(Get(server.URL).
		SelectorContains("#foo", "bar").
		Exec(cntx))
	a.Error(Get(server.URL).
		SelectorContains("#foo", "baXXX").
		Exec(cntx))
	a.NoError(Get(server.URL).
		SelectorContains("div.item", "title").
		Exec(cntx))
	a.Error(Get(server.URL).
		SelectorContains("", "nothing").
		Exec(cntx))
}

func Test_Http_Get500(t *testing.T) {
	a := assert.New(t)
	cntx := ExecutionContextImpl{}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(500)
		resp.Write([]byte(html))
	}))
	defer server.Close()

	a.Error(Get(server.URL).Exec(cntx))
	a.NoError(Get(server.URL).
		HasCode(500).
		Exec(cntx))
	a.NoError(Get(server.URL).
		HasCodeRange(400, 600).
		Exec(cntx))
}
