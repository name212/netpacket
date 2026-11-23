// Copyright 2026
// license that can be found in the LICENSE file.

package response

import (
	"fmt"
	"net/http"

	"github.com/name212/netpacket/l7/http/content"
	"github.com/name212/netpacket/l7/http/server/log"
	"github.com/name212/netpacket/l7/http/server/middleware"
)

const (
	NotFoundMsg = "not found"
	PanicMsg    = "internal server error"

	LogMsgErrWriteErrorResponse   = "cannot write error response"
	LogMsgErrWriteNotFondResponse = "cannot write not found response"
	LogMsgErrWritePanicResponse   = "cannot write panic response"
)

var (
	ErrNotFound = fmt.Errorf("%s", NotFoundMsg)
)

func Error(w http.ResponseWriter, r *http.Request, code int, msg string, err error) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	logger.RequestError(r, msg, err)

	l, errW := writeErrorContent(w, data, msg, err, code)
	if errW != nil {
		logger.RequestError(r, LogMsgErrWriteErrorResponse, errW)
	}

	logger.RequestErrorSent(r, code, l)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusNotFound, NotFoundMsg, ErrNotFound)
}

func InternalError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	Error(w, r, http.StatusInternalServerError, msg, err)
}

func Panic(w http.ResponseWriter, r *http.Request, err error) {
	data, r := middleware.GetServerData(r)

	_, errW := writeErrorContent(w, data, PanicMsg, err, http.StatusInternalServerError)

	if errW != nil {
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

func getErrorContent(w http.ResponseWriter, d *middleware.CtxServerData, msg string, err error, code int) []byte {
	out, ct := d.ErrorWriter()(msg, err)

	if ct == "" {
		ct = content.ContentText
	}

	w.Header().Add(content.ContentTypeHeader, ct)

	w.WriteHeader(code)

	return out
}

func writeErrorContent(w http.ResponseWriter, d *middleware.CtxServerData, msg string, err error, code int) (int, error) {
	out := getErrorContent(w, d, msg, err, code)
	return w.Write(out)
}
