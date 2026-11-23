// Copyright 2026
// license that can be found in the LICENSE file.

package content

import "net/http"

const ContentTypeHeader = "Content-Type"

type ContentType = string

const (
	ContentNo   ContentType = ""
	ContentText ContentType = "text/plain"
	ContentForm ContentType = "multipart/form-data"

	ContentJSON    ContentType = "application/json"
	ContentYAML    ContentType = "text/yaml"
	ContentAltYAML ContentType = "application/x-yaml"
)

func CustomContentType(t string) ContentType {
	return t
}

func GetResponseContentType(r *http.Response) ContentType {
	if r == nil || len(r.Header) == 0 {
		return ContentNo
	}

	return r.Header.Get(ContentTypeHeader)
}
