// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/name212/netpacket/utils/strings"
)

func SetCtxRequestID(next http.Handler, idLen ...int) http.Handler {
	l := 16
	if len(idLen) > 0 || idLen[0] > 0 {
		l = idLen[0]
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, setRequestID(r, l))
	})
}

func setRequestID(r *http.Request, l int) *http.Request {
	data, r := getAndSetServerData(r)
	if data.requestID != "" {
		return r
	}

	data.requestID = strings.RandString(l)
	return r
}
