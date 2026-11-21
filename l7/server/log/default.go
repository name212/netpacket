// Copyright 2026
// license that can be found in the LICENSE file.

package log

import "log/slog"

func ErrorDefault(msg string, err error) { 
	slog.Error(
		msg,
		errorAttrForLog(err),
	)
}

func PanicDefault(err error, trace string, bindAddr string) {
	slog.Error(
		PanicMsg,
		panicAttrForLog(err),
		panicTraceAttrForLog(trace),
		errorAttrForLog(err),
	)
}