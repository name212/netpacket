// Copyright 2026
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/name212/netpacket/utils/value"
)

var (
	loggersMu                    = &sync.RWMutex{}
	loggers   map[string]*Logger = make(map[string]*Logger)

	DefaultDestination io.Writer    = os.Stdout
	DefaultLevel       slog.Leveler = slog.LevelInfo
)

type InitParams struct {
	BindAddr            string
	RequestDataConsumer RequestDataConsumer
	NoAccessLogs        bool
}

func NewInitParams(bindAddr string, requestDataConsumer RequestDataConsumer) *InitParams {
	r := &InitParams{
		BindAddr:            bindAddr,
		RequestDataConsumer: requestDataConsumer,
		NoAccessLogs:        false,
	}

	r.validateOrPanic()

	return r
}

func InitLog(p *InitParams, parent ...*slog.Logger) *Logger {
	p.validateOrPanic()

	initDefaultIfNeed(p, parent)

	return setLogger(
		p.BindAddr,
		p,
		initLogger(p, parent),
	)
}

func initLogger(p *InitParams, parent []*slog.Logger) *slog.Logger {
	var res *slog.Logger

	if len(parent) > 0 && parent[0] != nil {
		h := parent[0].Handler()
		res = slog.New(h)
	} else {
		res = slog.New(slog.NewJSONHandler(DefaultDestination, &slog.HandlerOptions{
			Level: DefaultLevel,
		}))

		if slog.Default() == nil {
			slog.SetDefault(res)
		}
	}

	res = res.With(
		logtAttrForLog(),
		bindAddrAttrForLog(p.BindAddr),
	)

	return res
}

func getLogger(bindAddr string) *Logger {
	loggersMu.RLock()
	defer loggersMu.RUnlock()

	res, ok := loggers[bindAddr]
	if !ok {
		return nil
	}

	return res
}

func initDefaultIfNeed(p *InitParams, parent []*slog.Logger) {
	if def := getLogger(DefaultLoggerBind); def != nil {
		return
	}

	setLogger(DefaultLoggerBind, p, initLogger(p, parent))
}

func setLogger(bindAddr string, p *InitParams, logger *slog.Logger) *Logger {
	loggersMu.Lock()
	defer loggersMu.Unlock()

	res := &Logger{
		Logger:              logger,
		requestDataConsumer: p.RequestDataConsumer,
		noAccess:            p.NoAccessLogs,
	}

	loggers[bindAddr] = res

	return res
}

func (p *InitParams) validateOrPanic() {
	var msg []string

	if p != nil {
		if p.BindAddr == "" {
			msg = append(msg, "bind address empty")
		}

		if value.IsNil(p.RequestDataConsumer) {
			msg = append(msg, "; request data consumer is nil")
		}
	} else {
		msg = append(msg, "log init params is nil")
	}

	if len(msg) > 0 {
		panic(
			fmt.Sprintf("invalid init logger params: %s",
				strings.Join(msg, "; "),
			),
		)
	}
}

func (p *InitParams) WithNoAccessLogs() *InitParams {
	p.NoAccessLogs = true
	return p
}
