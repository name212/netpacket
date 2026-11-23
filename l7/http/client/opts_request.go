// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"io"
	"log/slog"
	"time"
)

type responseHeadersFilter func(k string) bool

type RequestOpts struct {
	*ClientOpts

	body                  io.ReadCloser
	responseHeadersDest   Headers
	responseHeadersFilter responseHeadersFilter
}

type RequestOpt func(o *RequestOpts)

func newRequestOpts() *RequestOpts {
	res := &RequestOpts{
		ClientOpts:            newClientOpts(),
		body:                  nil,
		responseHeadersDest:   nil,
		responseHeadersFilter: noResponseHeadersFilter,
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

func WithHeaders(h Headers) RequestOpt {
	return func(o *RequestOpts) {
		WithClientHeaders(h)(o.ClientOpts)
	}
}

func WithResponseHeaders(dest Headers, headersToConsume ...string) RequestOpt {
	return func(o *RequestOpts) {
		if dest == nil {
			return
		}

		o.responseHeadersDest = dest

		if len(headersToConsume) > 0 {
			o.responseHeadersFilter = getResponseHeadersFilter(headersToConsume)
		}
	}
}

func WithLogResponseBodyWhenError(l slog.Level, limit int64) RequestOpt {
	return func(o *RequestOpts) {
		WithClientLogResponseBodyWhenError(l, limit)(o.ClientOpts)
	}
}

func WithNoLogResponseBodyWhenError() RequestOpt {
	return func(o *RequestOpts) {
		WithClientNoLogResponseBodyWhenError()(o.ClientOpts)
	}
}

func WithHTTPS() RequestOpt {
	return func(o *RequestOpts) {
		WithClientHTTPS()(o.ClientOpts)
	}
}

func noResponseHeadersFilter(string) bool {
	return true
}

func getResponseHeadersFilter(headers []string) responseHeadersFilter {
	fm := make(map[string]struct{}, len(headers))
	for _, h := range headers {
		fm[h] = struct{}{}
	}

	return func(k string) bool {
		_, ok := fm[k]
		return ok
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
