// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
)

func SetCtxIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, setIP(r))
	})
}

func setIP(r *http.Request) *http.Request {
	ip := GetRealUserIPFromHeaders(r)
	data, r := getAndSetServerData(r)
	data.ip = ip
	return r
}

func GetRealUserIPFromHeaders(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")

	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}
