// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"

	"github.com/name212/netpacket/l7/server/log"
	"github.com/name212/netpacket/l7/server/middleware"
	"github.com/name212/netpacket/l7/server/response"
	"github.com/name212/netpacket/utils/value"
)

const (
	HealthZPath = "/_healthz"
	ReadyZPath  = "/_readyz"
)

var (
	okHealthZResponse = []byte("OK")
)

type (
	healthHandler  func(*http.Request) error
	HealthZHandler = healthHandler
	ReadyZHandler  = healthHandler
)

func getHandler(opts *middleWareParams) http.Handler {
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

	userHandler := opts.handler

	panicHandler := opts.panicHandler

	if value.IsNil(userHandler) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if panicHandler != nil {
				defer func() {
					err := recover()
					handlePanic(opts, r, err, panicHandler)
				}()
			}

			if r = setMiddlewares(w, r, opts); r == nil {
				return
			}

			if !healthPathsChecker(r.URL.Path) {
				response.WriteErrorResponse(w, r, http.StatusNotFound, "not found", nil)
				return
			}

			healthHandler.ServeHTTP(w, r)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if panicHandler != nil {
			defer func() {
				err := recover()
				handlePanic(opts, r, err, panicHandler)
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

	if !IsHead(r) || !IsGet(r) {
		response.WriteErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			fmt.Sprintf("incorrect method '%s' for healthz, should be '%s' or '%s'", r.Method, http.MethodHead, http.MethodGet),
			nil,
		)
		return
	}

	if err := h.doHealth(r); err != nil {
		response.WriteErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			"not health",
			err,
		)
	}

	if IsGet(r) {
		if _, err := w.Write(okHealthZResponse); err != nil {
			response.SendInternalError(w, r, "cannot write ok response for healthz handler", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
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

func handlePanic(opts *middleWareParams, r *http.Request, errAny any, handler PanicHandler) {
	if value.IsNil(errAny) {
		return
	}

	err := fmt.Errorf("panic: %v", errAny)

	opts.logger.RequestError(r, "server got panic", err)

	handler(errAny)
}
