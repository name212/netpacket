// Copyright 2026
// license that can be found in the LICENSE file.

package log

import (
	"log/slog"

	"github.com/name212/netpacket/l7/http/log"
)

const (
	SourceAddrAttr  = "client_ip"
	UserAttr        = log.UserAttr
	RequestIDAttr   = "request_id"
	MethodAttr      = log.MethodAttr
	PathAttr        = log.PathAttr
	UserAgentAttr   = log.UserAgentAttr
	ResponseLenAttr = log.ResponseLenAttr
	StatusCodeAttr  = log.StatusCodeAttr
	LogTypeAttr     = log.LogTypeAttr
	BindAddrAttr    = "bind_addr"

	ErrAttr = log.LogTypeAttr

	PanicAttr      = "panic"
	PanicTraceAttr = "trace"
)

const (
	LogTypeServerVal = "http-server"

	OKRequestMsg   = log.OKRequestMsg
	FailRequestMsg = log.FailRequestMsg

	DefaultLoggerBind = "unknown"

	PanicMsg = "server got panic"
)

func logtAttrForLog() LogAttr {
	return slog.String(LogTypeAttr, LogTypeServerVal)
}

func bindAddrAttrForLog(addr string) LogAttr {
	return slog.String(BindAddrAttr, addr)
}

func errorAttrForLog(err error) LogAttr {
	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}

	return slog.String(ErrAttr, errStr)
}

func panicAttrForLog(err error) LogAttr {
	return slog.Bool(PanicAttr, true)
}

func panicTraceAttrForLog(trace string) LogAttr {
	return slog.String(PanicTraceAttr, trace)
}
