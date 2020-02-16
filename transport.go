// Package logtransport implements logging wrapper for http.RoundTripper(a.k.a interface for http.Transport)
// As output it shows request and its headers as well as response code and response body if it was read in your application
package logtransport

import (
	"io"
	"net/http"
)

// Opts - logtransport options
type Opts struct {
	//LogMetadata - if set logs request's method, protocol and URI; response's status code and Content-Length. Default - true
	LogMetadata bool
	//LogReqHeaders - if set logs request's headers. Default - true
	LogReqHeaders bool
	//LogReqBody - if set logs request's body. Default - false
	LogReqBody bool
	//LogRespHeaders - if set logs response's headers. Default - true
	LogRespHeaders bool
	//LogRespBody - if set logs response's body. Only displayed is response is read by app. Default - false
	LogRespBody bool
}

var defaultOpts = &Opts{
	LogMetadata:    true,
	LogReqHeaders:  true,
	LogReqBody:     false,
	LogRespHeaders: true,
	LogRespBody:    false,
}

func New(tr http.RoundTripper, logOutput io.Writer, opts *Opts) http.RoundTripper {

	if opts == nil {
		opts = defaultOpts
	}
	if opts.LogMetadata {
		tr = &logRoundTripper{
			RoundTripper: tr,
			logOutput:    logOutput,
			opts:         opts,
		}
	}
	return tr
}
