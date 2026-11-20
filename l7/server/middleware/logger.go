// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/name212/netpacket/l7/server/log"
)

func SetCtxLogger(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, setLogger(r, logger))
	})
}

func setLogger(r *http.Request, logger *log.Logger) *http.Request {
	data, r := getAndSetServerData(r)
	data.logger = logger

	return r
}
