// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
)

type ErrorWriter func(msg string, err error) []byte

func TextErrorWriter(msg string, err error) []byte {
	return []byte(msg)
}

func SetCtxErrorWriter(next http.Handler, writer ErrorWriter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, setErrorWriter(r, writer))
	})
}

func setErrorWriter(r *http.Request, w ErrorWriter) *http.Request {
	if w == nil {
		w = TextErrorWriter
	}

	data, r := getAndSetServerData(r)
	data.errorWriter = w

	return r
}
