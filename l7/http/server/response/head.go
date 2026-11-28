// Copyright 2026
// license that can be found in the LICENSE file.

package response

import (
	"net/http"

	"github.com/name212/netpacket/l7/http/server/middleware"
)

func Head(w http.ResponseWriter, r *http.Request, code int, headers ...Headers) {
	writeHead(w, r, code, headers)
}

func writeHead(w http.ResponseWriter, r *http.Request, code int, headers []Headers) {
	_ = setHeaders(w, headers, -1)

	w.WriteHeader(code)

	data, r := middleware.GetServerData(r)

	data.Logger().RequestOKSent(r, code, 0)
}
