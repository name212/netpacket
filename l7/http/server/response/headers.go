// Copyright 2026
// license that can be found in the LICENSE file.

package response

import "net/http"

type Headers = map[string]string

func setHeaders(w http.ResponseWriter, headers []Headers) {
	if len(headers) == 0 {
		return
	}

	for _, hh := range headers {
		if len(hh) == 0 {
			continue
		}

		for k, v := range hh {
			w.Header().Add(k, v)
		}
	}
}

func getHeaders(headers []Headers) []Headers {
	hd := headers
	if len(hd) == 0 {
		hd = make([]Headers, 1)
	}
	if hd[0] == nil {
		hd[0] = make(Headers, 1)
	}

	return hd
}
