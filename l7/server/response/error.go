// Copyright 2026
// license that can be found in the LICENSE file.

package response

import (
	"fmt"
	"net/http"

	"github.com/name212/netpacket/l7/server/log"
	"github.com/name212/netpacket/l7/server/middleware"
)

const (
	NotFoundMsg = "not found"
	PanicMsg = "internal server error"

	LogMsgErrWriteErrorResponse = "cannot write error response"
	LogMsgErrWriteNotFondResponse = "cannot write not found response"
	LogMsgErrWritePanicResponse = "cannot write panic response"
)

var (
	ErrNotFound = fmt.Errorf("%s", NotFoundMsg)
)

func Error(w http.ResponseWriter, r *http.Request, code int, msg string, err error) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	logger.RequestError(r, msg, err)

	w.WriteHeader(code)

	out := data.ErrorWriter()(msg, err)
	if _, errW := w.Write(out); errW != nil {
		logger.RequestError(r, LogMsgErrWriteErrorResponse, errW)
	}

	logger.RequestErrorSent(r, code, len(out))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	code := http.StatusNotFound

	w.WriteHeader(code)

	out := data.ErrorWriter()(NotFoundMsg, ErrNotFound)
	if _, errW := w.Write(out); errW != nil {
		logger.RequestError(r, LogMsgErrWriteNotFondResponse, errW)
	}

	logger.RequestErrorSent(r, code, len(out))
}

func InternalError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	Error(w, r, http.StatusInternalServerError, msg, err)
}

func Panic(w http.ResponseWriter, r *http.Request, err error) {
	data, r := middleware.GetServerData(r)

	w.WriteHeader(http.StatusInternalServerError)

	out := data.ErrorWriter()(NotFoundMsg, ErrNotFound)
	if _, errW := w.Write(out); errW != nil {
		logger := data.Logger()

		if logger != nil {
			logger.RequestError(r, LogMsgErrWritePanicResponse, errW)
			return
		}
		log.ErrorDefault(LogMsgErrWritePanicResponse, errW)
	}
}

func LogInternalError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	logger.RequestError(r, msg, err)
}
