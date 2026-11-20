// Copyright 2026
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/name212/netpacket/l7/server/log"
)

var shutdownTimeout = 15 * time.Second

type server struct {
	httpServer *http.Server
	stopped    *atomic.Bool
	logger     *log.Logger
	addr       string
}

func newServer(bindAddr string, handler http.Handler, logger *log.Logger) *server {
	s := &http.Server{
		Addr:              bindAddr,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		MaxHeaderBytes:    16 * 1024,
		ErrorLog:          logger.ServerErrorLogger(),
		Handler:           handler,
	}

	return &server{
		httpServer: s,
		stopped:    &atomic.Bool{},
		logger:     logger,
		addr:       bindAddr,
	}
}

func (s *server) start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.shutdown(nil)
	}()

	s.logger.Info("Server started")

	var shutdownErr error
	if err := s.httpServer.ListenAndServe(); err != nil {
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if errStop := s.httpServer.Shutdown(shutdownCtx); errStop != nil {
		s.logger.Error("Cannot shutdown proxy HTTP-server gracefully", errStop)
	} else {
		s.logger.Info("HTTP-sever stopped")
	}

	return err
}
