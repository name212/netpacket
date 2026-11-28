// Copyright 2026
// license that can be found in the LICENSE file.

package content

import (
	"fmt"
	"net/http"
	"strings"
)

const ContentTypeHeader = "Content-Type"

type ContentType = string

const (
	TypeNo ContentType = ""

	TypeBytes ContentType = "application/octet-stream"
	TypeText  ContentType = "text/plain"
)

var textTypes = map[string]struct{}{
	TypeText:      {},
	TypeDocMD:     {},
	TypeScriptSh:  {},
	TypeScriptPHP: {},
	TypeJSON:      {},
	TypeYAML:      {},
	TypeYAMLAlt:   {},
	TypeXML:       {},
	TypeCSV:       {},
	TypeHTML:      {},
	TypeXHTML:     {},
	TypeCSS:       {},
	TypeJS:        {},
}

func CustomContentType(t string) ContentType {
	return t
}

func ResponseContentType(r *http.Response) ContentType {
	if r == nil || len(r.Header) == 0 {
		return TypeNo
	}

	ctFull := r.Header.Get(ContentTypeHeader)
	if ctFull == "" {
		return TypeNo
	}

	return strings.TrimSpace(strings.Split(ctFull, ";")[0])
}

func SetTypeHeader(w http.ResponseWriter, ct ContentType, encoding ...string) {
	w.Header().Add(ContentTypeHeader, getContentTypeHeader(ct, encoding))
}

func AppendContentType(dest map[string]string, ct ContentType, encoding ...string) {
	dest[ContentTypeHeader] = getContentTypeHeader(ct, encoding)
}

func IsText(t ContentType) bool {
	_, ok := textTypes[t]
	return ok
}

func getContentTypeHeader(ct ContentType, encoding []string) string {
	if ct == TypeNo {
		ct = TypeBytes
	}

	res := ct
	if IsText(ct) {
		enc := "utf-8"
		if len(encoding) > 0 && encoding[0] != "" {
			enc = encoding[0]
		}
		res = fmt.Sprintf("%s; charset=%s", res, enc)
	}

	return res
}
