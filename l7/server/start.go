// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/name212/netpacket/l7/server/log"
	"github.com/name212/netpacket/l7/server/middleware"
	"github.com/name212/netpacket/utils/value"
)

type BindAddr interface {
	Addr() string

	IsUnix() bool
	ExtractUnixSocketPath() (string, error)
	UnixSocketMode() (int, error)
}

func Start(ctx context.Context, bindAddr BindAddr, opts ...ServerOpt) error {
	if value.IsNil(bindAddr) {
		return fmt.Errorf("bind address is nil")
	}

	addr := bindAddr.Addr()
	if addr == "" {
		return fmt.Errorf("bind address is empty")
	}

	serverOpts := newServerOpts()

	for _, o := range opts {
		o(serverOpts)
	}

	loggerParams := log.NewInitParams(addr, logDataConsumer)
	if serverOpts.noAccessLogs {
		loggerParams.WithNoAccessLogs()
	}

	logger := log.InitLog(loggerParams, serverOpts.parentLogger)

	mParams := newMiddleWareParams(serverOpts, logger)

	handler := getHandler(mParams, addr)

	s, err := newServer(bindAddr, handler, logger, serverOpts)
	if err != nil {
		return err
	}

	return s.start(ctx)
}

func StartTopLevel(root context.Context, bindAddr BindAddr, opts ...ServerOpt) error {
	ctx, cancel := HandleTermSignalsCtx(root)

	defer cancel()

	return Start(ctx, bindAddr, opts...)
}

func HandleTermSignalsCtx(root context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		root,
		os.Interrupt,
		syscall.SIGTERM,
	)
}

func logDataConsumer(r *http.Request) log.RequestData {
	data, _ := middleware.GetServerData(r)
	return data
}

func appendErr(err error, f string, args ...any) error {
	curErr := fmt.Errorf(f, args...)
	if err != nil {
		return fmt.Errorf("%v\n%v", err, curErr)
	}

	return curErr
}
