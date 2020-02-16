package logtransport

import (
	"fmt"
	"io"
)

// Writes to LogOutput everything it reads from io.ReadCloser
type logReadCloser struct {
	io.Reader
	close        func() error
	logOutput    io.Writer
	notFirstRead bool
	prefix       string
}

func newReadCloser(prefix string, logOutput io.Writer, rc io.ReadCloser) io.ReadCloser {
	if rc == nil {
		return rc
	}

	return &logReadCloser{
		Reader:    io.TeeReader(rc, logOutput),
		close:     rc.Close,
		logOutput: logOutput,
		prefix:    prefix,
	}
}

func (rc *logReadCloser) Close() error {
	if rc.close == nil {
		return nil
	}
	return rc.close()
}

func (lr *logReadCloser) Read(p []byte) (n int, err error) {
	if !lr.notFirstRead {
		fmt.Fprint(lr.logOutput, lr.prefix)
		lr.notFirstRead = true
	}
	n, err = lr.Reader.Read(p)
	if err == io.EOF {
		fmt.Fprint(lr.logOutput, "\n\n")
	}
	return n, err

}
