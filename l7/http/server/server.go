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

	"github.com/name212/netpacket/l7/http/server/log"
	"github.com/name212/netpacket/utils/value"
)

var shutdownTimeout = 15 * time.Second

type server struct {
	httpServer *http.Server
	stopped    *atomic.Bool
	logger     *log.Logger
	l          listener
	shutdownCh chan struct{}
}

type listener interface {
	start() (net.Listener, error)
	shutdown() error
}

func newServer(bindAddr BindAddr, handler http.Handler, logger *log.Logger, opts *ServerOpts) (*server, error) {
	s := &http.Server{
		ReadTimeout:       opts.timeouts.read,
		ReadHeaderTimeout: opts.timeouts.readHeader,
		WriteTimeout:      opts.timeouts.write,

		Handler: handler,

		MaxHeaderBytes:               opts.headersOpts.maxHeaderBytes,
		MaxHeaderValueCount:          opts.headersOpts.maxHeaderValueCount,
		DisableGeneralOptionsHandler: opts.headersOpts.disableOptionsHandler,

		Protocols: opts.protocols,
	}

	if opts.timeouts.idle > 0 {
		s.IdleTimeout = opts.timeouts.idle
	}

	if opts.http2 != nil {
		if !opts.protocols.UnencryptedHTTP2() && opts.tls == nil {
			return nil, fmt.Errorf("tls config not passed when http2 enabled")
		}

		s.HTTP2 = opts.http2.cfg
	}

	if opts.tls != nil {
		s.TLSConfig = opts.tls
	}

	s.ErrorLog = logger.ServerErrorLogger()

	var l listener

	if bindAddr.IsUnix() {
		l = newUnixListener(bindAddr)
	} else {
		l = newPortListener(bindAddr)
	}

	return &server{
		httpServer: s,
		stopped:    &atomic.Bool{},
		logger:     logger,
		l:          l,
		shutdownCh: make(chan struct{}),
	}, nil
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

	shutdownErr = s.shutdown(shutdownErr)

	<-s.shutdownCh

	return shutdownErr
}

func (s *server) shutdown(err error) error {
	if !s.stopped.CompareAndSwap(false, true) {
		return err
	}

	defer func() {
		close(s.shutdownCh)
	}()

	var resErr = err

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	errStopHTTP := s.httpServer.Shutdown(shutdownCtx)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	addr BindAddr
}

func newPortListener(addr BindAddr) *portListener {
	return &portListener{
		addr: addr,
	}
}

func (p *portListener) start() (net.Listener, error) {
	addr := p.addr.Addr()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("cannot start port listener for '%s': %w", addr, err)
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
	addr BindAddr
	path string
}

func newUnixListener(bindAddr BindAddr) *unixListener {
	return &unixListener{
		addr: bindAddr,
	}
}

const maxUnixMode = 0o777

func (u *unixListener) start() (net.Listener, error) {
	mode := 0o700

	extractedMode, err := u.addr.UnixSocketMode()
	if err != nil {
		return nil, err
	}

	if extractedMode > 0 {
		if extractedMode > maxUnixMode {
			return nil, fmt.Errorf("incorrect unix socket mode %o should be <= %o in base 8", extractedMode, maxUnixMode)
		}
		mode = extractedMode
	}

	fMode := os.FileMode(mode)

	u.path, err = u.addr.ExtractUnixSocketPath()
	if err != nil {
		return nil, fmt.Errorf("cannot start listener with unix-socket '%s', cannot extract unix socket path: %w", u.l.Addr(), err)
	}

	if err := u.prepareSocksPath(); err != nil {
		return nil, err
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

func (u *unixListener) prepareSocksPath() error {
	statPath, err := os.Stat(u.path)
	if err == nil {
		if statPath.IsDir() {
			return fmt.Errorf(" unix socket path '%s' is dir", u.path)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot get stat for unix socket path '%s': %w", u.path, err)
	}

	socketDir := path.Dir(u.path)

	if socketDir == "" {
		return nil
	}

	dirStat, err := os.Stat(socketDir)
	if err == nil {
		if !dirStat.IsDir() {
			return fmt.Errorf("unix socket dir '%s' for path '%s' exists and not not dir", socketDir, u.path)
		}

		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot get stat for unix socket dir '%s': %w", socketDir, err)
	}

	if err := os.MkdirAll(socketDir, 0o777); err != nil {
		return fmt.Errorf(
			"cannot create socket dir '%s' for path '%s': %w",
			socketDir,
			u.path,
			err,
		)
	}

	return nil
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
