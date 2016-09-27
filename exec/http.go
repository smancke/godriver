package exec

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpExec struct {
	req              *http.Request
	initError        error
	expectations     []HttpExpectation
	codeExpectionSet bool
}

type HttpExpectation func(response *http.Response, body string) error

func Get(url string) *HttpExec {
	r, err := http.NewRequest("GET", url, nil)
	return &HttpExec{
		req:       r,
		initError: err,
	}
}

func Post(url string, contentType string, body string) *HttpExec {
	r, err := http.NewRequest("POST", url, strings.NewReader(body))
	r.Header.Set("Content-Type", contentType)
	return &HttpExec{
		req:       r,
		initError: err,
	}
}

func (httpExec *HttpExec) String() string {
	return fmt.Sprintf("->%v %v", httpExec.req.Method, httpExec.req.URL.String())
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
	httpExec.codeExpectionSet = true
	return httpExec
}

func (httpExec *HttpExec) HasCodeRange(min, max int) *HttpExec {
	httpExec.Expect(func(resp *http.Response, body string) error {
		if !(min <= resp.StatusCode && resp.StatusCode <= max) {
			return fmt.Errorf("response code was %v, but expected: %v <= code <= %v", resp, min, max)
		}
		return nil
	})
	httpExec.codeExpectionSet = true
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

func (httpExec *HttpExec) Exec(cntx ExecutionContext) error {
	if httpExec.initError != nil {
		return httpExec.initError
	}

	resp, err := http.DefaultClient.Do(httpExec.req)
	if err != nil {
		return err
	}
	if !httpExec.codeExpectionSet && resp.StatusCode != 200 {
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
