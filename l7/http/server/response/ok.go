// Copyright 2026
// license that can be found in the LICENSE file.

package response

import (
	"fmt"
	"net/http"

	"github.com/name212/netpacket/l7/http/content"
	"github.com/name212/netpacket/l7/http/server/middleware"
)

func OKWithStatus(w http.ResponseWriter, r *http.Request, code int, response []byte, headers ...Headers) {
	writeOK(w, r, code, response, headers)
}

func OK(w http.ResponseWriter, r *http.Request, response []byte, headers ...Headers) {
	writeOK(w, r, http.StatusOK, response, headers)
}

func Created(w http.ResponseWriter, r *http.Request, response []byte, headers ...Headers) {
	writeOK(w, r, http.StatusCreated, response, headers)
}

func Accepted(w http.ResponseWriter, r *http.Request, response []byte, headers ...Headers) {
	writeOK(w, r, http.StatusAccepted, response, headers)
}

func writeOK(w http.ResponseWriter, r *http.Request, code int, response []byte, headers []Headers) {
	data, r := middleware.GetServerData(r)

	logger := data.Logger()

	ct := setHeaders(w, headers)

	lenR := len(response)

	if lenR > 0 && ct == content.ContentNo {
		w.Header().Add(content.ContentTypeHeader, content.ContentText)
	}

	w.WriteHeader(code)

	if lenR > 0 {
		if _, errW := w.Write(response); errW != nil {
			logger.RequestError(
				r,
				fmt.Sprintf("cannot write ok response with code %d", code),
				errW,
			)
		}
	}

	logger.RequestOKSent(r, code, len(response))
}
