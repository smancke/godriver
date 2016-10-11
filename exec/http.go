package exec

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpExec struct {
	Method             string
	Url                string
	Header             http.Header
	Body               []byte
	expectations       []HttpExpectation
	codeExpectationSet bool
}

type HttpExpectation func(response *http.Response, body string) error

func Get(url string) *HttpExec {
	return &HttpExec{
		Method: "GET",
		Url:    url,
		Header: http.Header{},
	}
}

func Post(url string, contentType string, body string) *HttpExec {
	return &HttpExec{
		Method: "POST",
		Url:    url,
		Body:   []byte(body),
		Header: http.Header{"Content-Type": {contentType}},
	}
}

func (httpExec *HttpExec) WithAuthorization(authorizationHeader string) *HttpExec {
	httpExec.Header.Set("Authorization", authorizationHeader)
	return httpExec
}

func (httpExec *HttpExec) WithBasicAuth(username, password string) *HttpExec {
	enc := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return httpExec.WithAuthorization("Basic " + enc)
}

func (httpExec *HttpExec) String(cntx Context) string {
	return cntx.ExpandVarsNoError(fmt.Sprintf("->%v %v", httpExec.Method, httpExec.Url))
}

func (httpExec *HttpExec) Expect(e HttpExpectation) {
	if httpExec.expectations == nil {
		httpExec.expectations = make([]HttpExpectation, 0, 0)
	}
	httpExec.expectations = append(httpExec.expectations, e)
}

func (httpExec *HttpExec) HasContentType(contentType string) *HttpExec {
	httpExec.Expect(func(resp *http.Response, body string) error {
		if !strings.HasPrefix(resp.Header.Get("Content-Type"), contentType) {
			return fmt.Errorf("content type was %q, but expected: %q", resp.Header.Get("Content-Type"), contentType)
		}
		return nil
	})
	return httpExec
}

func (httpExec *HttpExec) HasNonErrorCode() *HttpExec {
	return httpExec.HasCodeRange(200, 399)
}

func (httpExec *HttpExec) HasCode(code int) *HttpExec {
	httpExec.Expect(func(resp *http.Response, body string) error {
		if resp.StatusCode != code {
			return fmt.Errorf("response code was %v, but expected: %v", resp, code)
		}
		return nil
	})
	httpExec.codeExpectationSet = true
	return httpExec
}

func (httpExec *HttpExec) HasCodeRange(min, max int) *HttpExec {
	httpExec.Expect(func(resp *http.Response, body string) error {
		if !(min <= resp.StatusCode && resp.StatusCode <= max) {
			return fmt.Errorf("response code was %v, but expected: %v <= code <= %v", resp, min, max)
		}
		return nil
	})
	httpExec.codeExpectationSet = true
	return httpExec
}

func (httpExec *HttpExec) Contains(substring string) *HttpExec {
	httpExec.Expect(func(resp *http.Response, body string) error {
		if !strings.Contains(body, substring) {
			return fmt.Errorf("response does not contain %q, but was: %q", substring, body)
		}
		return nil
	})
	return httpExec
}

func (httpExec *HttpExec) SelectorContains(selector, substring string) *HttpExec {
	httpExec.Expect(func(resp *http.Response, body string) error {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			return err
		}

		selection := doc.Find(selector)
		if selection.Length() == 0 {
			return fmt.Errorf("selection %q not found in response: %q", selector, body)
		}

		if !strings.Contains(selection.Text(), substring) {
			html, _ := selection.Html()
			return fmt.Errorf("selection does not contain %q, but was: %q", substring, html)
		}
		return nil
	})
	return httpExec
}

func (httpExec *HttpExec) Exec(cntx Context) error {
	url, err := cntx.ExpandVars(string(httpExec.Url))
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if len(httpExec.Body) > 0 {
		body, err := cntx.ExpandVars(string(httpExec.Body))
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(httpExec.Method, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header = http.Header{}
	for k := range httpExec.Header {
		v, err := cntx.ExpandVars(string(httpExec.Header.Get(k)))
		if err != nil {
			return err
		}
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !httpExec.codeExpectationSet && resp.StatusCode != 200 {
		return fmt.Errorf("response code was %v, but expected 200 (by default)", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	for _, expectation := range httpExec.expectations {
		err := expectation(resp, string(body))
		if err != nil {
			return err
		}
	}

	return nil
}
