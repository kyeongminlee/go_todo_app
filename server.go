package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func NewServer(listener net.Listener, mux http.Handler) *Server {
	return &Server{
		server:   &http.Server{Handler: mux},
		listener: listener,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.server.Serve(s.listener); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait()
}
