package configs

import (
	"context"
	"errors"
	"net/http"
	"time"

	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/models"

	"go.uber.org/fx"
)

type Server struct {
	port    int
	server  *http.Server
	handler http.Handler // ← Interface, not concrete type
	logger  logger.Logger
}

func NewServer(handler *http.Handler, cfg *models.ConfigEnv2) (*Server, error) {
	url, err := shared_kernel.ParseUrl(cfg.Server.HTTPListenUrl)
	if err != nil {
		return nil, errors.New("HTTPListenUrl is not set or is invalid; set it in config.yml")
	}
	return &Server{
		port: url.Port,
		server: &http.Server{
			Addr:         url.Address,
			Handler:      *handler,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}, nil
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func RegisterFX(lc fx.Lifecycle, server *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(); err != nil && err != http.ErrServerClosed {
					server.logger.Error("HTTP server error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
