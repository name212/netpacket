// Copyright 2026
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"log/slog"
	"net/http"
)

type (
	RequestData interface {
		IP() string
		UserID() string
		RequestID() string
	}

	RequestDataConsumer func(*http.Request) RequestData

	LogAttr = slog.Attr

	Logger struct {
		*slog.Logger
		requestDataConsumer RequestDataConsumer
		noAccess bool
	}

	logFunc func(string, ...any)
)

func (l *Logger) ServerErrorLogger() *log.Logger {
	return slog.NewLogLogger(l.Logger.Handler(), slog.LevelError)
}

func (l *Logger) Str(key, val string) LogAttr {
	return slog.String(key, val)
}

func (l *Logger) Info(msg string, attrbs ...any) {
	if l == nil {
		return
	}

	l.Logger.Info(msg, attrbs...)
}

func (l *Logger) Error(msg string, err error) {
	if l == nil {
		return
	}

	l.Logger.Error(msg, errorAttrForLog(err))
}

func (l *Logger) RequestAttr(level slog.Level, r *http.Request, msg string, attrs ...any) {
	if l == nil {
		return
	}

	l.logRequestAttr(l.getLoggerFunc(level), r, msg, attrs)
}

func (l *Logger) Request(level slog.Level, r *http.Request, msg string) {
	if l == nil {
		return
	}

	l.logRequest(l.getLoggerFunc(level), r, msg)
}

func (l *Logger) RequestDebug(r *http.Request, msg string) {
	if l == nil {
		return
	}

	l.logRequest(l.Logger.Debug, r, msg)
}

func (l *Logger) RequestInfo(r *http.Request, msg string) {
	if l == nil {
		return
	}

	l.logRequest(l.Logger.Info, r, msg)
}

func (l *Logger) RequestWarn(r *http.Request, msg string) {
	if l == nil {
		return
	}

	l.logRequest(l.Logger.Warn, r, msg)
}

func (l *Logger) RequestError(r *http.Request, msg string, err error) {
	if l == nil {
		return
	}

	data := l.requestDataConsumer(r)

	l.Logger.Error(
		msg,
		slog.String(SourceAddrAttr, data.IP()),
		slog.String(UserAttr, data.UserID()),
		slog.String(RequestIDAttr, data.RequestID()),
		slog.String(MethodAttr, r.Method),
		slog.String(PathAttr, r.URL.Path),
		slog.String(UserAgentAttr, r.UserAgent()),
		errorAttrForLog(err),
	)
}

func (l *Logger) RequestPanic(r *http.Request, err error, trace string) {
	if l == nil {
		return
	}

	data := l.requestDataConsumer(r)

	l.Logger.Error(
		PanicMsg,
		slog.String(SourceAddrAttr, data.IP()),
		slog.String(UserAttr, data.UserID()),
		slog.String(RequestIDAttr, data.RequestID()),
		slog.String(MethodAttr, r.Method),
		slog.String(PathAttr, r.URL.Path),
		slog.String(UserAgentAttr, r.UserAgent()),
		panicAttrForLog(err),
		panicTraceAttrForLog(trace),
		errorAttrForLog(err),
	)
}

func (l *Logger) RequestErrorSent(r *http.Request, code, responseLen int) {
	if l == nil {
		return
	}

	if l.noAccess {
		return
	}

	data := l.requestDataConsumer(r)

	l.Logger.Error(
		FailRequestMsg,
		slog.String(SourceAddrAttr, data.IP()),
		slog.String(UserAttr, data.UserID()),
		slog.String(RequestIDAttr, data.RequestID()),
		slog.String(MethodAttr, r.Method),
		slog.String(PathAttr, r.URL.Path),
		slog.String(UserAgentAttr, r.UserAgent()),
		slog.Int(ResponseLenAttr, responseLen),
		slog.Int(StatusCodeAttr, code),
	)
}

func (l *Logger) RequestOKSent(r *http.Request, code, responseLen int) {
	if l == nil {
		return
	}

	if l.noAccess {
		return
	}

	data := l.requestDataConsumer(r)

	l.Logger.Info(
		OKRequestMsg,
		slog.String(SourceAddrAttr, data.IP()),
		slog.String(UserAttr, data.UserID()),
		slog.String(RequestIDAttr, data.RequestID()),
		slog.String(MethodAttr, r.Method),
		slog.String(PathAttr, r.URL.Path),
		slog.String(UserAgentAttr, r.UserAgent()),
		slog.Int(ResponseLenAttr, responseLen),
		slog.Int(StatusCodeAttr, code),
	)
}

func (l *Logger) getLoggerFunc(level slog.Level) logFunc {
	var logger logFunc

	switch level {
	case slog.LevelInfo:
		logger = l.Logger.Info
	case slog.LevelDebug:
		logger = l.Logger.Debug
	case slog.LevelWarn:
		logger = l.Logger.Warn
	case slog.LevelError:
		logger = l.Logger.Error
	default:
		logger = l.Logger.Info
	}

	return logger
}

func (l *Logger) logRequestAttr(logger logFunc, r *http.Request, msg string, attrs []any) {
	if l == nil {
		return
	}

	data := l.requestDataConsumer(r)

	if len(attrs) > 0 {
		attrs = append(
			attrs,
			slog.String(SourceAddrAttr, data.IP()),
			slog.String(UserAttr, data.UserID()),
			slog.String(RequestIDAttr, data.RequestID()),
			slog.String(MethodAttr, r.Method),
			slog.String(PathAttr, r.URL.Path),
			slog.String(UserAgentAttr, r.UserAgent()),
		)

		logger(
			msg,
			attrs...,
		)

		return
	}

	l.logRequest(logger, r, msg)
}

func (l *Logger) logRequest(logger logFunc, r *http.Request, msg string) {
	data := l.requestDataConsumer(r)

	logger(
		msg,
		slog.String(SourceAddrAttr, data.IP()),
		slog.String(UserAttr, data.UserID()),
		slog.String(RequestIDAttr, data.RequestID()),
		slog.String(MethodAttr, r.Method),
		slog.String(PathAttr, r.URL.Path),
		slog.String(UserAgentAttr, r.UserAgent()),
	)
}
