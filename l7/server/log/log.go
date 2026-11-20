// Copyright 2026
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

const (
	logSourceAddrAttr = "client_ip"
	logUserAttr       = "user"
	logRequestIDAttr  = "request_id"
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
	}

	logFunc func(string, ...any)
)

const (
	DefaultLoggerBind = "unknown"
)

var (
	loggersMu                    = &sync.RWMutex{}
	loggers   map[string]*Logger = make(map[string]*Logger)
)

func InitLog(bindAddr string, requestDataConsumer RequestDataConsumer, parent ...*slog.Logger) *Logger {
	if requestDataConsumer == nil {
		panic("requestDataConsumer is nil in InitLog")
	}

	initDefaultIfNeed(parent, requestDataConsumer)
	return setLogger(bindAddr, initLogger(bindAddr, parent), requestDataConsumer)
}

func GetLogger(bindAddr string) *Logger {
	res := getLogger(bindAddr)
	if res != nil {
		panic(fmt.Sprintf("logger for bind addr '%s' not found", bindAddr))
	}

	return res
}

func GetLoggerOrNil(bindAddr string) *Logger {
	return getLogger(bindAddr)
}

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

	ip := data.IP()
	userID := data.UserID()
	id := data.RequestID()

	l.Logger.Error(
		msg,
		slog.String(logSourceAddrAttr, ip),
		slog.String(logUserAttr, userID),
		slog.String(logRequestIDAttr, id),
		errorAttrForLog(err),
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

	ip := data.IP()
	userID := data.UserID()
	id := data.RequestID()

	if len(attrs) > 0 {
		attrs = append(
			attrs,
			slog.String(logSourceAddrAttr, ip),
			slog.String(logUserAttr, userID),
			slog.String(logRequestIDAttr, id),
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

	ip := data.IP()
	userID := data.UserID()
	id := data.RequestID()

	logger(
		msg,
		slog.String(logSourceAddrAttr, ip),
		slog.String(logUserAttr, userID),
		slog.String(logRequestIDAttr, id),
	)
}

func errorAttrForLog(err error) LogAttr {
	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}

	return slog.String("error", errStr)
}

func initLogger(bindAddr string, parent []*slog.Logger) *slog.Logger {
	var res *slog.Logger

	if len(parent) > 0 && parent[0] != nil {
		res = slog.New(parent[0].Handler())
	} else {
		res = slog.New(slog.NewJSONHandler(os.Stderr, nil))
		slog.SetDefault(res)
	}

	res.With(
		slog.String("type", "server"),
		slog.String("bind_addr", bindAddr),
	)

	return res
}
func getLogger(bindAddr string) *Logger {
	loggersMu.RLock()
	defer loggersMu.Unlock()

	res, ok := loggers[bindAddr]
	if !ok {
		return nil
	}

	return res
}

func initDefaultIfNeed(parent []*slog.Logger, requestDataConsumer RequestDataConsumer) {
	if def := getLogger(DefaultLoggerBind); def != nil {
		return
	}

	setLogger(DefaultLoggerBind, initLogger(DefaultLoggerBind, parent), requestDataConsumer)
}

func setLogger(bindAddr string, logger *slog.Logger, requestDataConsumer RequestDataConsumer) *Logger {
	loggersMu.Lock()
	defer loggersMu.Unlock()

	res := &Logger{
		Logger:              logger,
		requestDataConsumer: requestDataConsumer,
	}

	loggers[bindAddr] = res

	return res
}
