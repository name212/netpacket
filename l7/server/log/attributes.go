// Copyright 2026
// license that can be found in the LICENSE file.

package log

import "log/slog"

const (
	SourceAddrAttr  = "client_ip"
	UserAttr        = "user"
	RequestIDAttr   = "request_id"
	MethodAttr      = "method"
	PathAttr        = "path"
	UserAgentAttr   = "agent"
	ResponseLenAttr = "response_len"
	StatusCodeAttr  = "status"
	LogTypeAttr     = "logt"
	BindAddrAttr    = "bind_addr"

	ErrAttr = "error"

	PanicAttr      = "panic"
	PanicTraceAttr = "trace"
)

const (
	LogTypeServerVal = "server"

	OKRequestMsg   = "Request OK"
	FailRequestMsg = "Request error"

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
