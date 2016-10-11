package exec

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	cntx := &ContextImpl{}

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(html))
	}))
	defer server.Close()

	a.Error(Get("h :// invalid").Exec(cntx))
	a.NoError(Get(server.URL).Exec(cntx))

	// Content-type
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
	cntx := &ContextImpl{}

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

func Test_Http_Post(t *testing.T) {
	a := assert.New(t)
	cntx := NewDefaultContext()

	server := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// echo service
		resp.Header().Set("Content-Type", req.Header.Get("Content-Type"))
		body, _ := ioutil.ReadAll(req.Body)
		resp.Write(body)
	}))
	defer server.Close()

	a.NoError(Post(server.URL, "application/foo", "demo data").
		HasContentType("application/foo").
		Contains("demo data").
		Exec(cntx))

	a.Error(Post(server.URL, "application/foo", "demo data").
		HasContentType("application/bar").
		Exec(cntx))

	a.Error(Post(server.URL, "application/foo", "demo data").
		Contains("wrong data").
		Exec(cntx))

	a.Error(Post("h :// invalid", "application/foo", "demo data").
		Exec(cntx))
}
