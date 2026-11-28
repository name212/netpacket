// Copyright 2026
// license that can be found in the LICENSE file.

package response

import (
	"net/http"

	"github.com/name212/netpacket/l7/http/content"
)

type Headers = map[string]string

func setHeaders(w http.ResponseWriter, headers []Headers, contentLen int) content.ContentType {
	if len(headers) == 0 {
		if contentLen >= 0 {
			content.SetLenHeader(w, contentLen)
		}
		return content.TypeNo
	}

	ct := content.TypeNo

	hasLenHeader := false

	for _, hh := range headers {
		if len(hh) == 0 {
			continue
		}

		for k, v := range hh {
			if k == content.ContentTypeHeader {
				ct = content.CustomContentType(v)
			}

			if k == content.LenHeader {
				hasLenHeader = true
			}

			w.Header().Add(k, v)
		}
	}

	if !hasLenHeader && contentLen >= 0 {
		content.SetLenHeader(w, contentLen)
	}

	return ct
}

func HeadersWithContentType(ct content.ContentType, headers ...Headers) Headers {
	hd := getHeaders(headers)[0]
	content.AppendContentType(hd, ct)
	return hd
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
