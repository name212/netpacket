// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/name212/netpacket/l7/http/log"
)

const (
	LogHostAttr         = "host"
	LogIsUnixSocketAttr = "unix_socket"
	LogBodyAttr         = "body"
	LogMethodAttr       = log.MethodAttr
	LogPathAttr         = log.PathAttr
	LogUserAttr         = log.UserAttr
	LogUserAgentAttr    = log.UserAgentAttr
	LogRequestLenAttr   = log.RequestLenAttr
	LogResponseLenAttr  = log.ResponseLenAttr
	LogStatusCodeAttr   = log.StatusCodeAttr
	LogTypeAttr         = log.LogTypeAttr

	LogErrAttr = log.ErrAttr
)

const (
	LogTypeClientVal = "http-client"

	LogOKRequestMsg    = "Request OK"
	LogFailRequestMsg  = "Request error"
	LogResponseBodyMsg = "response body"
)

type logFunc func(string, ...any)
type logFuncConsumer func(l *slog.Logger) logFunc

func logFuncErr(l *slog.Logger) logFunc {
	return l.Error
}

func logFuncInfo(l *slog.Logger) logFunc {
	return l.Info
}

func (c *Client) logAttrs(requestOpts *RequestOpts, r *http.Request, code int, attrs ...any) []any {
	a := []any{
		slog.String(LogTypeAttr, LogTypeClientVal),
		slog.String(LogMethodAttr, r.Method),
		slog.String(LogPathAttr, r.URL.Path),
		slog.String(LogHostAttr, c.host),
	}

	if id := requestOpts.getUserID(); id != "" {
		a = append(a, slog.String(LogUserAttr, id))
	}

	if agent := requestOpts.getUserAgent(); agent != "" {
		a = append(a, slog.String(LogUserAgentAttr, agent))
	}

	if code > 0 {
		a = append(a, slog.Int(LogStatusCodeAttr, code))
	}

	if c.isUnix {
		a = append(a, slog.Bool(LogIsUnixSocketAttr, true))
	}

	a = append(a, attrs...)

	return a
}

func (c *Client) getLogger(requestOpts *RequestOpts) *slog.Logger {
	logger := requestOpts.logger
	if logger == nil {
		logger = c.logger
	}

	return logger
}

func (c *Client) readBodyLogger(requestOpts *RequestOpts, r *http.Request, code int) *responseBodyReadLogger {
	logger := c.getLogger(requestOpts)
	if logger == nil {
		return nil
	}

	return newResponseBodyReadLogger(
		logger,
		c.logAttrs(requestOpts, r, code),
	)
}

func (c *Client) logBody(requestOpts *RequestOpts, respBody *ResponseBody, r *http.Request) {
	logger := c.getLogger(requestOpts)
	if logger == nil {
		return
	}

	lf := requestOpts.responseBodyErrLogFunc(logger)

	if lf == nil {
		return
	}

	limit := requestOpts.logBodyLimit

	buf := bytes.NewBuffer(nil)

	n, err := respBody.CopyKeep(buf, limit)

	bodyStr := "<no body>"

	if n > 0 {
		bodyStr = buf.String()
	} else if err != nil {
		bodyStr = fmt.Sprintf("<read body error: %v>", err)
	}

	attrs := []any{
		slog.String(LogBodyAttr, bodyStr),
	}

	if n <= limit {
		attrs = append(attrs, slog.Int64(LogResponseLenAttr, n))
	}

	attrs = c.logAttrs(requestOpts, r, 0, attrs...)

	lf(LogResponseBodyMsg, attrs...)
}

func (c *Client) log(requestOpts *RequestOpts, msg string, fc logFuncConsumer, r *http.Request, code int, attrs ...any) {
	logger := c.getLogger(requestOpts)
	if logger == nil {
		return
	}

	fc(logger)(msg, c.logAttrs(requestOpts, r, code, attrs...)...)
}

type responseBodyReadLogger struct {
	logger *slog.Logger
	attrs  []any

	bodyLogger logFuncConsumer
	bodyLimit  int64

	lenAttrAdded bool
}

func newResponseBodyReadLogger(logger *slog.Logger, attrs []any) *responseBodyReadLogger {
	return &responseBodyReadLogger{
		logger: logger,
		attrs:  attrs,
	}
}

func (l *responseBodyReadLogger) read(n int64, sufMsg string) {
	if l == nil || l.logger == nil {
		return
	}

	la := len(l.attrs)

	if !l.lenAttrAdded {
		a := make([]any, 0, la+1)
		if la > 0 {
			a = append(a, l.attrs...)
		}

		a = append(a, slog.Int64(LogResponseLenAttr, n))

		l.attrs = a
		l.lenAttrAdded = true
	}

	l.attrs[la-1] = slog.Int64(LogResponseLenAttr, n)

	l.logger.Debug(
		fmt.Sprintf("read response body %s bytes %d", sufMsg, n),
		l.attrs...,
	)
}

// func (l *responseBodyReadLogger) debug(msg string) {
// 	if l == nil || l.logger == nil {
// 		return
// 	}

// 	l.logger.Debug(
// 		msg,
// 		l.attrs...,
// 	)
// }

func (l *responseBodyReadLogger) body(b []byte, msg string) {
	if l == nil || l.logger == nil {
		return
	}

	if l.bodyLogger == nil {
		return
	}

	lf := l.bodyLogger(l.logger)
	if lf == nil {
		return
	}

	bLen := int64(len(b))

	if bLen < 1 {
		return
	}

	logLen := l.bodyLimit

	if logLen > bLen {
		logLen = bLen
	}

	a := make([]any, 0, len(l.attrs)+1)
	a = append(a, l.attrs...)
	a = append(a, slog.String(LogBodyAttr, string(b[:logLen])))

	lf(msg, a...)
}
