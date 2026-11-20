// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	htls "crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/name212/netpacket/l7/server/middleware"
	"github.com/name212/netpacket/utils/value"
)

type PanicHandler func(any)

func DummyPanicHandler(any) {}

type (
	Timeouts struct {
		read       time.Duration
		readHeader time.Duration
		write      time.Duration
		idle       time.Duration
	}

	HeadersOpts struct {
		disableOptionsHandler bool
		maxHeaderBytes        int
		maxHeaderValueCount   int
	}

	HTTP2Opts struct {
		cfg         *http.HTTP2Config
		unencrypted bool
		only        bool
	}

	ServerOpts struct {
		errorWriter    middleware.ErrorWriter
		handler        http.Handler
		userIDConsumer middleware.UserIDConsumer
		requestIDLen   int
		panicHandler   PanicHandler
		parentLogger   *slog.Logger
		healthZ        HealthZHandler
		readyZ         ReadyZHandler
		timeouts       *Timeouts
		headersOpts    *HeadersOpts
		tls            *htls.Config
		http2          *HTTP2Opts
		protocols      *http.Protocols
	}

	ServerOpt func(*ServerOpts)
)

func newServerOpts() *ServerOpts {
	p := &http.Protocols{}
	p.SetHTTP1(true)
	p.SetHTTP2(true)
	p.SetUnencryptedHTTP2(false)

	return &ServerOpts{
		errorWriter:    middleware.TextErrorWriter,
		handler:        nil,
		userIDConsumer: nil,
		requestIDLen:   16,
		panicHandler:   nil,
		parentLogger:   nil,
		timeouts:       NewTimeouts(),
		headersOpts:    NewHeadersOpts(),
		tls:            nil,
		http2:          nil,
		protocols:      p,
	}
}

func WithErrorWriter(errorWriter middleware.ErrorWriter) ServerOpt {
	return func(so *ServerOpts) {
		if errorWriter != nil {
			so.errorWriter = errorWriter
		}
	}
}

func WithHandler(handler http.Handler) ServerOpt {
	return func(so *ServerOpts) {
		if !value.IsNil(handler) {
			so.handler = handler
		}
	}
}

func WithUserIDConsumer(c middleware.UserIDConsumer) ServerOpt {
	return func(so *ServerOpts) {
		so.userIDConsumer = c
	}
}

func WithDefaultUserIDConsumer() ServerOpt {
	return func(so *ServerOpts) {
		so.userIDConsumer = middleware.ReadUserNameFromHeader
	}
}

func WithRequestIDLen(l int) ServerOpt {
	return func(so *ServerOpts) {
		if l > 0 {
			so.requestIDLen = l
		}
	}
}

func WithPanicHandler(h PanicHandler) ServerOpt {
	return func(so *ServerOpts) {
		so.panicHandler = h
	}
}

func WithParentLogger(logger *slog.Logger) ServerOpt {
	return func(so *ServerOpts) {
		so.parentLogger = logger
	}
}

func WithHealthZHandler(h HealthZHandler) ServerOpt {
	return func(so *ServerOpts) {
		so.healthZ = h
	}
}

func WithReadyZHandler(h ReadyZHandler) ServerOpt {
	return func(so *ServerOpts) {
		so.readyZ = h
	}
}

func WithTimeouts(t *Timeouts) ServerOpt {
	return func(so *ServerOpts) {
		if t != nil {
			so.timeouts = t
		}
	}
}

func WithHeadersOpts(h *HeadersOpts) ServerOpt {
	return func(so *ServerOpts) {
		if h != nil {
			so.headersOpts = h
		}
	}
}

func WithTLS(t *htls.Config) ServerOpt {
	return func(so *ServerOpts) {
		if t != nil {
			so.tls = t
		}
	}
}

func WithTLSKeys(certPath string, keyPath string) (ServerOpt, error) {
	certs, err := htls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot load server certificates from cert path '%s' and key path '%s': %w",
			certPath,
			keyPath,
			err,
		)
	}

	cfg := &htls.Config{
		Certificates: []htls.Certificate{certs},
	}

	return func(so *ServerOpts) {
		so.tls = cfg
	}, nil
}

func WithDisableHTTP2() ServerOpt {
	return func(so *ServerOpts) {
		so.protocols.SetHTTP1(true)
		so.protocols.SetHTTP2(false)
		so.protocols.SetUnencryptedHTTP2(false)
	}
}

func WithHTTP2(opts *HTTP2Opts) ServerOpt {
	return func(so *ServerOpts) {
		if opts == nil {
			return
		}

		so.http2 = opts

		so.protocols.SetHTTP2(true)

		if opts.unencrypted {
			so.protocols.SetUnencryptedHTTP2(true)
		}

		if opts.only {
			so.protocols.SetHTTP1(false)
		}
	}
}

var DefaultTimeout = 15 * time.Second

func NewTimeouts() *Timeouts {
	return &Timeouts{
		read:       DefaultTimeout,
		readHeader: DefaultTimeout,
		write:      DefaultTimeout,
	}
}

func (t *Timeouts) WithRead(d time.Duration) *Timeouts {
	t.read = d
	return t
}

func (t *Timeouts) WithReadHeader(d time.Duration) *Timeouts {
	t.readHeader = d
	return t
}

func (t *Timeouts) WithWrite(d time.Duration) *Timeouts {
	t.write = d
	return t
}

func (t *Timeouts) WithIdle(d time.Duration) *Timeouts {
	t.idle = d
	return t
}

func NewHeadersOpts() *HeadersOpts {
	return &HeadersOpts{
		disableOptionsHandler: true,
		maxHeaderBytes:        16 * 1024,
		maxHeaderValueCount:   0,
	}
}

func (h *HeadersOpts) WithDisableOptionsHandler(d bool) *HeadersOpts {
	h.disableOptionsHandler = d
	return h
}

func (h *HeadersOpts) WithMaxHeaderBytes(b int) *HeadersOpts {
	if b >= 0 {
		h.maxHeaderBytes = b
	}
	return h
}

func (h *HeadersOpts) WithMaxHeaderCount(b int) *HeadersOpts {
	if b >= 0 {
		h.maxHeaderValueCount = b
	}
	return h
}

func (h *HeadersOpts) InfinityHeaders() *HeadersOpts {
	return h.WithMaxHeaderBytes(0).WithMaxHeaderCount(0)
}

func NewHTTP2Opts(cfg *http.HTTP2Config) *HTTP2Opts {
	if cfg == nil {
		panic("pass nil config to NewHTTP2Opts")
	}

	return &HTTP2Opts{
		cfg:         cfg,
		unencrypted: false,
		only:        true,
	}
}

func (h *HTTP2Opts) EnableUnencrypted() *HTTP2Opts {
	h.unencrypted = true
	return h
}

func (h *HTTP2Opts) OnlyHTTP2() *HTTP2Opts {
	h.only = true
	return h
}
