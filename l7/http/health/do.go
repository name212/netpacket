// Copyright 2026
// license that can be found in the LICENSE file.

package health

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/name212/netpacket/l7/http/client"
	"github.com/name212/netpacket/l7/http/log"
	"github.com/name212/netpacket/l7/http/server"
	"github.com/name212/netpacket/l7/params/bind"
)

var userAgentOpt = client.WithUserAgent("health-checker/0.0")
var timeoutOpt = client.WithTimeout(3 * time.Second)

func doHealth(ctx context.Context, addr *bind.Addr, path string, logger []*slog.Logger) error {
	l, msg := getLogger(path, logger)

	err := client.Head(ctx, addr, path, userAgentOpt, timeoutOpt, client.WithLogger(l))
	if err != nil {
		return fmt.Errorf("server '%s' is not %s: %w", addr.Addr(), msg, err)
	}

	return nil
}

func doHealthAsCmd(p *bind.AddrParams, cmdName, path string, logger []*slog.Logger) {
	if len(os.Args) > 1 && os.Args[1] != cmdName {
		return
	}

	f, err := bind.InitBindFlags(p)
	if err != nil {
		logAndExitError(path, err, logger)
		return
	}

	flag.Parse()

	addr, err := bind.GetBindAddress(f)
	if err != nil {
		logAndExitError(path, err, logger)
		return
	}

	ctx, cancel := server.HandleTermSignalsCtx(context.Background())
	defer cancel()

	if err := doHealth(ctx, addr, path, logger); err != nil {
		logAndExitError(path, err, logger)
		return
	}

	logAndExitOK(path, logger)
}

func logAndExitError(path string, err error, logger []*slog.Logger) {
	l, msg := getLogger(path, logger)

	if l != nil {
		l.Error(
			fmt.Sprintf("%s command failed", msg),
			slog.String(log.ErrAttr, err.Error()),
		)
	} else {
		fmt.Fprintf(os.Stderr, "%s command failed: %v\n", msg, err)
	}

	os.Exit(1)
}

func logAndExitOK(path string, logger []*slog.Logger) {
	l, msg := getLogger(path, logger)

	if l != nil {
		l.Info(msg)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}

	os.Exit(0)
}

func getLogger(path string, logger []*slog.Logger) (*slog.Logger, string) {
	var l *slog.Logger
	if len(logger) > 0 && logger[0] != nil {
		l = logger[0]
	}

	return l, getMsg(path)
}

func getMsg(path string) string {
	msg := "health"
	if path == server.ReadyZPath {
		msg = "ready"
	}

	return msg
}
