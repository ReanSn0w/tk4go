package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

// NewServer creates a new server instance
func NewServer(logger tools.Logger) *Server {
	return &Server{
		log: logger,
	}
}

type (
	Server struct {
		log tools.Logger
		srv *http.Server
	}
)

// Run starts the server
func (s *Server) Run(ctx context.Context, port int, h http.Handler) {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: s.prepareHandler(h),
	}

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.log.Logf("[ERROR] server failed: %v", err)
		} else {
			s.log.Logf("[INFO] server stopped")
		}

		cancel := tools.GetGlobalCancel(ctx)
		cancel(err)
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Logf("[ERROR] server shutdown failed: %v", err)
	} else {
		s.log.Logf("[INFO] server shutdown")
	}

	return err
}

func (s *Server) prepareHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
			return
		}

		h.ServeHTTP(w, r)
	})
}
