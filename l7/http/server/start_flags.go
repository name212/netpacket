// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"sync"

	"github.com/name212/netpacket/l7/params/bind"
)

type StartAllServers struct {
	BindAddr BindAddr
	Opts     []ServerOpt
}

type StartAllServerFlags struct {
	BindFlags *bind.AddrFlags
	Opts      []ServerOpt
}

type StartAllServerParams struct {
	BindFlags *bind.AddrParams
	Opts      []ServerOpt
}

func StartFromParsedFlags(root context.Context, bindFlags *bind.AddrFlags, opts ...ServerOpt) error {
	bindAddr, err := bind.GetBindAddress(bindFlags)
	if err != nil {
		return err
	}

	return Start(root, bindAddr, opts...)
}

func StartTopLevelFromParsedFlags(root context.Context, bindFlags *bind.AddrFlags, opts ...ServerOpt) error {
	ctx, cancel := HandleTermSignalsCtx(root)

	defer cancel()

	return StartFromParsedFlags(ctx, bindFlags, opts...)
}

func StartTopLevelFromFlagsSet(root context.Context, bindParams *bind.AddrParams, set *flag.FlagSet, opts ...ServerOpt) error {
	bindFlags, err := bind.InitBindFlags(bindParams, set)
	if err != nil {
		return err
	}

	flag.Parse()

	return StartTopLevelFromParsedFlags(root, bindFlags, opts...)
}

func StartTopLevelFromFlags(root context.Context, bindParams *bind.AddrParams, opts ...ServerOpt) error {
	return StartTopLevelFromFlagsSet(root, bindParams, nil, opts...)
}

func StartAll(root context.Context, wgLogger *slog.Logger, servers ...StartAllServers) error {
	serversCount := len(servers)

	if serversCount == 0 {
		return fmt.Errorf("cannot pass any servers to StartAll")
	}

	type result struct {
		bindAddr string
		err      error
	}

	ctx, cancel := context.WithCancel(root)

	canceler := newCanceler(cancel)

	resultsCh := make(chan *result, serversCount)

	for _, s := range servers {
		go func() {
			err := Start(ctx, s.BindAddr, s.Opts...)
			if err != nil {
				wgLogger.Error("Server returns error. Stop all...", slog.String("bind_addr", s.BindAddr.Addr()))
				canceler.do()
			}

			resultsCh <- &result{
				bindAddr: s.BindAddr.Addr(),
				err:      err,
			}
		}()
	}

	results := make([]*result, 0, serversCount)

	for r := range resultsCh {
		results = append(results, r)
		if len(results) == serversCount {
			break
		}
	}

	close(resultsCh)

	var serversErr error

	for _, r := range results {
		if r.err == nil {
			continue
		}

		curErr := fmt.Errorf("  bind address %s -> %v", r.bindAddr, r.err)
		if serversErr == nil {
			serversErr = curErr
			continue
		}

		serversErr = fmt.Errorf("%v\n%v", serversErr, curErr)
	}

	if serversErr != nil {
		return fmt.Errorf("some servers returns error:\n%v", serversErr)
	}

	return nil
}

func StartAllTopLevel(root context.Context, wgLogger *slog.Logger, servers ...StartAllServers) error {
	ctx, cancel := HandleTermSignalsCtx(root)
	defer cancel()

	return StartAll(ctx, wgLogger, servers...)
}

func StartAllFromFlagsTopLevel(root context.Context, wgLogger *slog.Logger, servers ...StartAllServerFlags) error {
	ctx, cancel := HandleTermSignalsCtx(root)
	defer cancel()

	return StartAllFromFlags(ctx, wgLogger, servers...)
}

func StartAllFromFlags(root context.Context, wgLogger *slog.Logger, servers ...StartAllServerFlags) error {
	serversCount := len(servers)

	if serversCount == 0 {
		return fmt.Errorf("cannot pass any servers to StartAllFromFlags")
	}

	var flagsErr error

	binds := make([]StartAllServers, 0, serversCount)

	for i, s := range servers {
		bindAddr, err := bind.GetBindAddress(s.BindFlags)
		if err != nil {
			flagsErr = appendErr(flagsErr, "  cannot extract bind addr for server %d: %v", i, err)
			continue
		}

		binds = append(binds, StartAllServers{
			BindAddr: bindAddr,
			Opts:     s.Opts,
		})
	}

	if flagsErr != nil {
		return fmt.Errorf("bind addresses extract errors:\n%v", flagsErr)
	}

	return StartAll(root, wgLogger, binds...)
}

func StartAllFromParams(root context.Context, wgLogger *slog.Logger, servers ...StartAllServerParams) error {
	serversCount := len(servers)

	if serversCount == 0 {
		return fmt.Errorf("cannot pass any servers to StartAllFromParams")
	}

	var flagsErr error

	flags := make([]StartAllServerFlags, 0, serversCount)

	for i, s := range servers {
		f, err := bind.InitBindFlags(s.BindFlags)
		if err != nil {
			flagsErr = appendErr(flagsErr, "  cannot extract bind addr for server %d: %v", i, err)
			continue
		}

		flags = append(flags, StartAllServerFlags{
			BindFlags: f,
			Opts:      s.Opts,
		})
	}

	if flagsErr != nil {
		return fmt.Errorf("bind addresses init flags errors:\n%v", flagsErr)
	}

	flag.Parse()

	return StartAllFromFlags(root, wgLogger, flags...)
}

func StartAllFromParamsTopLevel(root context.Context, wgLogger *slog.Logger, servers ...StartAllServerParams) error {
	ctx, cancel := HandleTermSignalsCtx(root)
	defer cancel()

	return StartAllFromParams(ctx, wgLogger, servers...)
}

type canceler struct {
	mu     sync.Mutex
	cancel context.CancelFunc
}

func newCanceler(cancel context.CancelFunc) *canceler {
	return &canceler{
		cancel: cancel,
	}
}

func (c *canceler) do() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel == nil {
		return
	}

	c.cancel()
	c.cancel = nil
}
