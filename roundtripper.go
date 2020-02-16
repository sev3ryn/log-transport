package logtransport

import (
	"fmt"
	"io"
	"net/http"
)

type logRoundTripper struct {
	http.RoundTripper
	logOutput io.Writer
	opts      *Opts
}

func (l *logRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	fmt.Fprintln(l.logOutput, ">", req.Method, req.Proto, req.RequestURI)
	if l.opts.LogReqHeaders {
		for k, v := range req.Header {
			fmt.Fprintln(l.logOutput, ">", k+":", v)
		}
	}

	if l.opts.LogReqBody {
		req.Body = newReadCloser("> Request Body:\n\n", l.logOutput, req.Body)
	}

	resp, err := l.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	fmt.Fprintln(l.logOutput, "<", resp.Proto, resp.Status)

	if l.opts.LogRespHeaders {
		for k, v := range resp.Header {
			fmt.Fprintln(l.logOutput, "<", k+":", v)
		}
	}

	if l.opts.LogRespBody {
		resp.Body = newReadCloser("< Response Body:\n\n", l.logOutput, resp.Body)
	}

	return resp, nil
}
