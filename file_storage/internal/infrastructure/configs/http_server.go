package configs

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"

	"go.uber.org/fx"
)

type Server struct {
	port    int
	server  *http.Server
	handler http.Handler // ← Interface, not concrete type
	logger  logger.Logger
}

func NewServer(handler *http.Handler, cfg *models.ConfigEnv) *Server {
	port, _ := strconv.Atoi(cfg.FileStorageHttpPort)
	return &Server{
		port: port,
		server: &http.Server{
			Addr:         cfg.FileStorageHttpAddr,
			Handler:      *handler,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
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
