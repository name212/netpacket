// Copyright 2026
// license that can be found in the LICENSE file.

package health

import (
	"context"
	"log/slog"

	"github.com/name212/netpacket/l7/http/server"
	"github.com/name212/netpacket/l7/params/bind"
)

const ReadyCmd = "_ready"

func DoReadyZ(ctx context.Context, addr *bind.Addr, logger ...*slog.Logger) error {
	return doHealth(ctx, addr, server.ReadyZPath, logger)
}

func DoReadyZFromFlags(ctx context.Context, f *bind.AddrFlags, logger ...*slog.Logger) error {
	addr, err := bind.GetBindAddress(f)
	if err != nil {
		return err
	}

	return DoHealthZ(ctx, addr, logger...)
}

func DoReadyZAsCmd(p *bind.AddrParams, cmdName string, logger ...*slog.Logger) {
	doHealthAsCmd(p, cmdName, server.ReadyZPath, logger)
}

func DoReadyZCmd(p *bind.AddrParams, logger ...*slog.Logger) {
	DoReadyZAsCmd(p, ReadyCmd, logger...)
}
