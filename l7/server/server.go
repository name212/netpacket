// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"sync/atomic"
	"time"

	"github.com/name212/netpacket/l7/params/bind"
	"github.com/name212/netpacket/l7/server/log"
	"github.com/name212/netpacket/utils/value"
)

var shutdownTimeout = 15 * time.Second

type server struct {
	httpServer *http.Server
	stopped    *atomic.Bool
	logger     *log.Logger
	l          listener
}

type listener interface {
	start() (net.Listener, error)
	shutdown() error
}

func newServer(bindAddr *bind.Addr, handler http.Handler, logger *log.Logger) *server {
	s := &http.Server{
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		MaxHeaderBytes:    16 * 1024,
		ErrorLog:          logger.ServerErrorLogger(),
		Handler:           handler,
	}

	var l listener

	if bindAddr.IsUnix() {
		l = newUnixListener(*bindAddr)
	} else {
		l = newPortListener(bindAddr)
	}

	return &server{
		httpServer: s,
		stopped:    &atomic.Bool{},
		logger:     logger,
		l:          l,
	}
}

func (s *server) start(ctx context.Context) error {
	listener, err := s.l.start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		_ = s.shutdown(nil)
	}()

	s.logger.Info("Server started")

	var shutdownErr error
	if err := s.httpServer.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			shutdownErr = fmt.Errorf("server closed with error: %v", err)
		}
	}

	return s.shutdown(shutdownErr)
}

func (s *server) shutdown(err error) error {
	if !s.stopped.CompareAndSwap(false, true) {
		return err
	}

	var resErr = err

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if errStopHTTP := s.httpServer.Shutdown(shutdownCtx); errStopHTTP != nil {
		s.logger.Error("Cannot shutdown HTTP-server gracefully. HTTP-server not stopped", errStopHTTP)
		resErr = appendErr(err, "http server not stopped skip stop listener")
	} else {
		if errStopListener := s.l.shutdown(); errStopListener != nil {
			s.logger.Error("Cannot shutdown proxy HTTP-server gracefully. Listener not stopped", errStopListener)
			resErr = appendErr(err, "http server listener not stopped")
		}

		s.logger.Info("HTTP-sever stopped")
	}

	return resErr
}

type portListener struct {
	l    net.Listener
	addr *bind.Addr
}

func newPortListener(addr *bind.Addr) *portListener {
	return &portListener{
		addr: addr,
	}
}

func (p *portListener) start() (net.Listener, error) {
	l, err := net.Listen("tcp", p.addr.BindAddr)
	if err != nil {
		return nil, fmt.Errorf("cannot start port listener: %w", err)
	}

	p.l = l
	return p.l, nil
}

func (p *portListener) shutdown() error {
	if err := stopListener(p.l); err != nil {
		return fmt.Errorf("cannot stop port listener: %w", err)
	}

	p.l = nil
	return nil
}

type unixListener struct {
	l    net.Listener
	addr *bind.Addr
	path string
}

func newUnixListener(addr bind.Addr) *unixListener {
	return &unixListener{
		addr: &addr,
	}
}

func (u *unixListener) start() (net.Listener, error) {
	mode := 0o777

	var err error

	if u.addr.UnixSocketMode > 0 {
		var mValidated int
		mValidated, err := bind.ValidateUnixSocketMode(fmt.Sprintf("%d", u.addr.UnixSocketMode))
		if err != nil {
			return nil, fmt.Errorf("cannot start listener with unix-socket '%s', incorrect socket permission mode: %w", u.addr.BindAddr, err)
		}
		mode = mValidated
	}

	fMode := os.FileMode(mode)

	u.path, err = u.addr.ExtractUnixSocketPath()
	if err != nil {
		return nil, fmt.Errorf("cannot start listener with unix-socket '%s', cannot extract unix socket path: %w", u.addr.BindAddr, err)
	}

	socketDir := path.Dir(u.path)

	if socketDir != "" {
		if err := os.MkdirAll(socketDir, 0o777); err != nil {
			return nil, fmt.Errorf(
				"cannot start listener with unix-socket '%s', cannot create socket dir '%s': %w",
				u.addr.BindAddr,
				socketDir,
				err,
			)
		}
	}

	l, err := net.Listen("unix", u.path)
	if err != nil {
		return nil, fmt.Errorf("cannot start listener with unix-socket '%s': %w", u.path, err)
	}

	u.l = l

	if err := os.Chmod(u.path, fMode); err != nil {
		stopErrStr := ""
		if stopErr := u.shutdown(); stopErr != nil {
			stopErrStr = fmt.Sprintf("; listener not stopped: %v", stopErr)
		}

		return nil, fmt.Errorf("cannot chmod to %o socket file '%s': %w%s", fMode, u.path, err, stopErrStr)
	}

	u.l = l
	return u.l, nil
}

func (u *unixListener) shutdown() error {
	if err := stopListener(u.l); err != nil {
		return fmt.Errorf("cannot stop unix listener '%s': %w file not removed", u.path, err)
	}

	if err := os.Remove(u.path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("unix listener '%s' stopped, but file not removed: %w", u.path, err)
		}
	}

	u.l = nil
	return nil
}

func stopListener(l net.Listener) error {
	if value.IsNil(l) {
		return nil
	}

	if err := l.Close(); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			return err
		}
	}

	return nil
}
