// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/name212/netpacket/l7/http/server/log"
	"github.com/name212/netpacket/l7/http/server/middleware"
	"github.com/name212/netpacket/l7/http/server/response"
	"github.com/name212/netpacket/utils/value"
)

const (
	HealthZPath = "/_healthz"
	ReadyZPath  = "/_readyz"

	OKHealthZResponse = "OK"
)

var (
	okHealthZResponse = []byte(OKHealthZResponse)
)

type (
	healthHandler  func(*http.Request) error
	HealthZHandler = healthHandler
	ReadyZHandler  = healthHandler
)

func getHandler(opts *middleWareParams, bindAddr string) http.Handler {
	healthHandler := &healthCheckHandler{
		healthZ: opts.healthZ,
		readyZ:  opts.readyZ,
	}

	healthPathsChecker := func(p string) bool {
		return p == HealthZPath
	}

	if opts.readyZ != nil {
		healthPathsChecker = func(p string) bool {
			return p == HealthZPath || p == ReadyZPath
		}
	}

	panicHandler := opts.panicHandler

	handlePanic := func(w http.ResponseWriter, r *http.Request, errAny any, trace []byte) {
		err := fmt.Errorf("panic: %v", errAny)

		t := string(trace)

		if opts.logger != nil {
			opts.logger.RequestPanic(r, err, t)
		} else {
			log.PanicDefault(err, t, bindAddr)
		}

		response.Panic(w, r, err)

		panicHandler(errAny, t)
	}

	userHandler := opts.handler

	if value.IsNil(userHandler) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if panicHandler != nil {
				defer func() {
					err := recover()
					if value.IsNil(err) {
						return
					}
					handlePanic(w, r, err, debug.Stack())
				}()
			}

			if r = setMiddlewares(w, r, opts); r == nil {
				return
			}

			if !healthPathsChecker(r.URL.Path) {
				response.NotFound(w, r)
				return
			}

			healthHandler.ServeHTTP(w, r)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if panicHandler != nil {
			defer func() {
				err := recover()
				if value.IsNil(err) {
					return
				}
				handlePanic(w, r, err, debug.Stack())
			}()
		}

		if r = setMiddlewares(w, r, opts); r == nil {
			return
		}

		if healthPathsChecker(r.URL.Path) {
			healthHandler.ServeHTTP(w, r)
			return
		}

		userHandler.ServeHTTP(w, r)
	})
}

func setMiddlewares(w http.ResponseWriter, r *http.Request, opts *middleWareParams) *http.Request {
	if opts.userIDConsumer != nil {
		return middleware.SetAllGeneralCtxAllParamsOnly(w, r, opts)
	}

	return middleware.SetGeneralCtxParamsOnly(r, opts)
}

type healthCheckHandler struct {
	readyZ  ReadyZHandler
	healthZ HealthZHandler
}

func (h *healthCheckHandler) doHealth(r *http.Request) error {
	var handler healthHandler

	switch r.URL.Path {
	case HealthZPath:
		handler = h.healthZ
	case ReadyZPath:
		handler = h.readyZ
	default:
		return nil
	}

	if handler != nil {
		if err := handler(r); err != nil {
			return fmt.Errorf("%s user handler returns error: %w", r.URL.Path, err)
		}
	}

	return nil
}

func (h *healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer BodyClose(r)

	if !IsHead(r) && !IsGet(r) {
		response.Error(
			w,
			r,
			http.StatusBadRequest,
			fmt.Sprintf("incorrect method '%s' for healthz, should be '%s' or '%s'", r.Method, http.MethodHead, http.MethodGet),
			nil,
		)
		return
	}

	if err := h.doHealth(r); err != nil {
		response.Error(
			w,
			r,
			http.StatusBadRequest,
			"not health",
			err,
		)
	}

	if IsGet(r) {
		response.OK(w, r, okHealthZResponse)
		return
	}

	response.Head(w, r, http.StatusOK)
}

type middleWareParams struct {
	*ServerOpts
	logger *log.Logger
}

func newMiddleWareParams(opts *ServerOpts, logger *log.Logger) *middleWareParams {
	return &middleWareParams{
		ServerOpts: opts,
		logger:     logger,
	}
}

func (p *middleWareParams) ErrorWriter() middleware.ErrorWriter {
	return p.errorWriter
}

func (p *middleWareParams) Logger() *log.Logger {
	return p.logger
}

func (p *middleWareParams) UserIDConsumer() middleware.UserIDConsumer {
	return p.userIDConsumer
}

func (p *middleWareParams) RequestIDLen() int {
	return p.requestIDLen
}
