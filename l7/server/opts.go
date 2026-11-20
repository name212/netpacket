// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"log/slog"
	"net/http"

	"github.com/name212/netpacket/l7/server/middleware"
	"github.com/name212/netpacket/utils/value"
)

type PanicHandler func(any)

type (
	ServerOpts struct {
		errorWriter    middleware.ErrorWriter
		handler        http.Handler
		userIDConsumer middleware.UserIDConsumer
		requestIDLen   int
		panicHandler   PanicHandler
		parentLogger   *slog.Logger
		healthZ        HealthZHandler
		readyZ         ReadyZHandler
	}

	ServerOpt func(*ServerOpts)
)

func newServerOpts() *ServerOpts {
	return &ServerOpts{
		errorWriter:    middleware.TextErrorWriter,
		handler:        nil,
		userIDConsumer: nil,
		requestIDLen:   16,
		panicHandler:   nil,
		parentLogger:   nil,
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
