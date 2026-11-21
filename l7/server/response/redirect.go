// Copyright 2026
// license that can be found in the LICENSE file.

package response

import "net/http"

func Redirect(w http.ResponseWriter, r *http.Request, url string, headers ...Headers) {
	writeRedirect(w, r, http.StatusFound, url, headers)
}

func RedirectTemp(w http.ResponseWriter, r *http.Request, url string, headers ...Headers) {
	writeRedirect(w, r, http.StatusTemporaryRedirect, url, headers)
}

func Moved(w http.ResponseWriter, r *http.Request, url string, headers ...Headers) {
	writeRedirect(w, r, http.StatusMovedPermanently, url, headers)
}

func writeRedirect(w http.ResponseWriter, r *http.Request, code int, url string, headers []Headers) {
	hd := getHeaders(headers)

	hd[0]["Location"] = url

	writeHead(w, r, code, hd)
}
