// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"log/slog"
	"time"

	"github.com/name212/netpacket/l7/http/server/middleware"
)

const (
	DefaultTimeoutSeconds = 10
	userAgentHeader       = "User-Agent"
)

type Headers = map[string]string

type ClientOpts struct {
	headers               Headers
	timeout               time.Duration
	schema                string
	logger                *slog.Logger
	noReplaceHost         bool
	urlPreparator         URLPreparator
	statusChecker         StatusChecker
	autoCloseResponseBody bool
}

type ClientOpt func(o *ClientOpts)

func newClientOpts() *ClientOpts {
	return &ClientOpts{
		headers:               nil,
		timeout:               DefaultTimeoutSeconds * time.Second,
		schema:                "http",
		noReplaceHost:         false,
		urlPreparator:         DummyURLPreparator,
		statusChecker:         StatusChecker200X,
		autoCloseResponseBody: true,
	}
}

func WithClientTimeout(timeout time.Duration) ClientOpt {
	return func(o *ClientOpts) {
		if timeout > 0 {
			o.timeout = timeout
		}
	}
}

func WithClientLogger(logger *slog.Logger) ClientOpt {
	return func(o *ClientOpts) {
		if logger != nil {
			o.logger = logger
		}
	}
}

func WithClientNoReplaceHost() ClientOpt {
	return func(o *ClientOpts) {
		o.noReplaceHost = true
	}
}

func WithClientStatusChecker(c StatusChecker) ClientOpt {
	return func(o *ClientOpts) {
		if c != nil {
			o.statusChecker = c
		}
	}
}

func WithClientURLPreparator(p URLPreparator) ClientOpt {
	return func(o *ClientOpts) {
		if p != nil {
			o.urlPreparator = p
		}
	}
}

func WithClientUserID(id string) ClientOpt {
	return func(o *ClientOpts) {
		if id != "" {
			o.addHeader(middleware.UserIDHeader, id)
		}
	}
}

func WithClientUserAgent(a string) ClientOpt {
	return func(o *ClientOpts) {
		if a != "" {
			o.addHeader(userAgentHeader, a)
		}
	}
}

func WithClientHTTPS() ClientOpt {
	return func(o *ClientOpts) {
		o.schema = "https"
	}
}

func WithClientNoAutoCloseResponseBody() ClientOpt {
	return func(o *ClientOpts) {
		o.autoCloseResponseBody = false
	}
}

func StatusChecker200X(code int) bool {
	return code >= 200 && code < 299
}

func (o *ClientOpts) prepareHeaders() {
	if o.headers == nil {
		o.headers = make(Headers)
	}
}

func (o *ClientOpts) addHeader(k, v string) {
	o.prepareHeaders()
	o.headers[k] = v
}

func (o *ClientOpts) addHeaders(h Headers) {
	o.prepareHeaders()
	for k, v := range h {
		o.headers[k] = v
	}
}

func (o *ClientOpts) getHeader(name string) string {
	res := ""

	if len(o.headers) > 0 {
		if h, ok := o.headers[name]; ok {
			res = h
		}

	}

	return res
}

func (o *ClientOpts) getUserID() string {
	return o.getHeader(middleware.UserIDHeader)
}

func (o *ClientOpts) getUserAgent() string {
	return o.getHeader(userAgentHeader)
}
