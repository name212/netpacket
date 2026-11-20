// Copyright 2026
// license that can be found in the LICENSE file.

package response

import (
	"net/http"

	"github.com/name212/netpacket/l7/server/middleware"
)

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, code int, msg string, err error) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	logger.RequestError(r, msg, err)

	out := data.ErrorWriter()(msg, err)
	if _, errW := w.Write(out); errW != nil {
		logger.RequestError(r, "cannot write error response", err)
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)
}

func WriteInternalErrorResponse(w http.ResponseWriter, r *http.Request, msg string, err error) {
	WriteErrorResponse(w, r, http.StatusInternalServerError, msg, err)
}

func SendInternalError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	logger.RequestError(r, msg, err)

	w.WriteHeader(http.StatusInternalServerError)
}
