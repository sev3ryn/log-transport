// Package logtransport implements logging wrapper for http.RoundTripper(a.k.a interface for http.Transport)
// As output it shows request and its headers as well as response code and response body if it was read in your application
package logtransport

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type headers map[string]string

func newRequest(method string, body io.Reader, h headers) *http.Request {
	req, _ := http.NewRequest(method, "http://dummy.com", body)
	req.Header.Set("Host", "example.com")
	for k, v := range h {
		req.Header.Set(k, v)
	}
	return req
}

func TestNew(t *testing.T) {
	type args struct {
		respCode    int
		respBody    io.Reader
		respHeaders headers
		req         *http.Request
	}

	tests := []struct {
		name          string
		args          *args
		wantLogOutput string
	}{
		{
			name: "metadata test",
			args: &args{
				req:         newRequest(http.MethodGet, strings.NewReader("hi server!"), headers{"from": "client"}),
				respCode:    http.StatusAccepted,
				respBody:    strings.NewReader("hello client!"),
				respHeaders: headers{"from": "server"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, v := range tt.args.respHeaders {
					w.Header().Add(k, v)
				}
				io.Copy(ioutil.Discard, r.Body)
				defer r.Body.Close()
				w.WriteHeader(tt.args.respCode)
				io.Copy(w, tt.args.respBody)
			}))
			defer ts.Close()

			logOutput := &bytes.Buffer{}
			cl := http.Client{
				Transport: New(http.DefaultTransport, logOutput, nil),
			}

			req, _ := http.NewRequest(http.MethodGet, ts.URL, tt.args.req.Body)
			req.Header = tt.args.req.Header
			resp, err := cl.Do(req)
			if err != nil {
				t.Fatal("no response from server:", err)
			}
			defer resp.Body.Close()

			io.Copy(ioutil.Discard, resp.Body)

			// not finished yet
			t.Logf("\n%s", logOutput.String())
		})
	}
}
