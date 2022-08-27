package internalhttp

import (
	"context"
	"net/http"

	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	srv  *http.Server
	app  app.Application
	logg Logger
}

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
}

func NewServer(logger Logger, app app.Application, addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(logger, http.HandlerFunc(EmptyHandler)))
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		app:  app,
		logg: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Debug("Server starting ...")

	if err := s.srv.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Debug("Server stopping ...")

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func EmptyHandler(w http.ResponseWriter, req *http.Request) {
}
