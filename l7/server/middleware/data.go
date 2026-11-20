// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"

	"github.com/name212/netpacket/l7/server/log"
)

type ctxDataKeyT string

const ctxDataKey ctxDataKeyT = "server-data"
const UnknownVal = "unknown"

type CtxServerData struct {
	ip          string
	userID      string
	requestID   string
	errorWriter ErrorWriter
	logger      *log.Logger
}

func GetServerData(r *http.Request) (*CtxServerData, *http.Request) {
	anyData := r.Context().Value(ctxDataKey)
	if data, ok := anyData.(*CtxServerData); ok {
		return data, r
	}

	return nil, r
}

func getAndSetServerData(r *http.Request) (*CtxServerData, *http.Request) {
	var data *CtxServerData
	data, r= GetServerData(r)
	if data != nil {
		return data, r
	}

	data = &CtxServerData{}

	r = r.WithContext(
		context.WithValue(r.Context(), ctxDataKey, data),
	)

	return data, r
}

func (d *CtxServerData) IP() string {
	if d == nil {
		return UnknownVal
	}

	return d.ip
}

func (d *CtxServerData) UserID() string {
	if d == nil {
		return UnknownVal
	}

	return d.userID
}

func (d *CtxServerData) ErrorWriter() ErrorWriter {
	res := TextErrorWriter

	if d != nil && d.errorWriter != nil {
		res = d.errorWriter
	}

	return res
}

func (d *CtxServerData) Logger() *log.Logger {
	if d == nil || d.logger == nil {
		return log.GetLogger(log.DefaultLoggerBind)
	}

	return d.logger
}

func (d *CtxServerData) RequestID() string {
	if d == nil || d.requestID == "" {
		return UnknownVal
	}

	return d.requestID
}
