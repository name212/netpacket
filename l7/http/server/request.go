// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/name212/netpacket/l7/http/server/middleware"
	"github.com/name212/netpacket/utils/value"
)

func BodyClose(r *http.Request) {
	if r.Body == nil || value.IsNil(r.Body) {
		return
	}

	if err := r.Body.Close(); err != nil {
		data, _ := middleware.GetServerData(r)
		data.Logger().RequestError(r, "cannot close request body", err)
	}
}

func IsHead(r *http.Request) bool {
	return r.Method == http.MethodHead
}

func IsGet(r *http.Request) bool {
	return r.Method == http.MethodGet
}
