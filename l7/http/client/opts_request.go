// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"io"
	"log/slog"
	"time"
)

type RequestOpts struct {
	*ClientOpts
	body io.ReadCloser
}

type RequestOpt func(o *RequestOpts)

func newRequestOpts() *RequestOpts {
	res := &RequestOpts{
		ClientOpts: newClientOpts(),
		body:       nil,
	}

	return res
}

func WithTimeout(timeout time.Duration) RequestOpt {
	return func(o *RequestOpts) {
		WithClientTimeout(timeout)(o.ClientOpts)
	}
}

func WithStatusChecker(c StatusChecker) RequestOpt {
	return func(o *RequestOpts) {
		WithClientStatusChecker(c)(o.ClientOpts)
	}
}

func WithLogger(logger *slog.Logger) RequestOpt {
	return func(o *RequestOpts) {
		WithClientLogger(logger)(o.ClientOpts)
	}
}

func WithBodyBytes(b []byte) RequestOpt {
	return func(o *RequestOpts) {
		if len(b) == 0 {
			return
		}

		o.body = newReaderCloseBytes(b)
	}
}

func WithBody(b string) RequestOpt {
	return func(o *RequestOpts) {
		if b == "" {
			return
		}

		WithBodyBytes([]byte(b))
	}
}

func WithBodyReader(b io.Reader) RequestOpt {
	return func(o *RequestOpts) {
		o.body = newReaderClose(b)
	}
}

func WithBodyReadCloser(b io.ReadCloser) RequestOpt {
	return func(o *RequestOpts) {
		o.body = newReaderClose(b)
	}
}

func WithNoAutoCloseResponseBody() RequestOpt {
	return func(o *RequestOpts) {
		o.autoCloseResponseBody = false
	}
}

func WithAutoCloseResponseBody() RequestOpt {
	return func(o *RequestOpts) {
		o.autoCloseResponseBody = true
	}
}

func WithUserID(id string) RequestOpt {
	return func(o *RequestOpts) {
		WithClientUserID(id)(o.ClientOpts)
	}
}

func WithUserAgent(a string) RequestOpt {
	return func(o *RequestOpts) {
		WithClientUserAgent(a)(o.ClientOpts)
	}
}

func WithHTTPS() RequestOpt {
	return func(o *RequestOpts) {
		WithClientHTTPS()(o.ClientOpts)
	}
}

func clientsOptsToRequest(opts []ClientOpt) []RequestOpt {
	if len(opts) == 0 {
		return nil
	}

	res := make([]RequestOpt, 0, len(opts))

	for _, opt := range opts {
		f := func(o *RequestOpts) {
			opt(o.ClientOpts)
		}

		res = append(res, f)
	}

	return res
}
