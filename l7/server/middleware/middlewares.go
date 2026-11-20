// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/name212/netpacket/l7/server/log"
)

type GeneralParams interface {
	ErrorWriter() ErrorWriter
	Logger() *log.Logger
}

type UserParams interface {
	UserIDConsumer() UserIDConsumer
	RequestIDLen() int
}

type AllParams interface {
	GeneralParams
	UserParams
}

func SetGeneralCtxParamsOnly(r *http.Request, params GeneralParams) *http.Request {
	return setErrorWriter(setLogger(r, params.Logger()), params.ErrorWriter())
}

func SetAllGeneralCtxAllParamsOnly(w http.ResponseWriter, r *http.Request, params AllParams) *http.Request {
	r = SetGeneralCtxParamsOnly(r, params)
	return setUserID(
		w, 
		setIP(
			setRequestID(r, params.RequestIDLen()),
		), 
		params.UserIDConsumer(), 
	)
}

func SetGeneralCtxParams(r *http.Request, params GeneralParams, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r = SetGeneralCtxParamsOnly(r, params); r != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func SetAllGeneralCtxParams(w http.ResponseWriter, r *http.Request, params AllParams, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r = SetAllGeneralCtxAllParamsOnly(w, r, params); r != nil {
			next.ServeHTTP(w, r)
		}
	})
}
