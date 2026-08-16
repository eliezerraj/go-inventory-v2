package webserver

import (
	"context"
	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"

	"github.com/go-inventory-v2/application/infrastructure/application"
	"github.com/go-inventory-v2/cmd/webserver/framework/fiber"
	"github.com/go-inventory-v2/application/config"

	//"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	cfg	*config.Config
	fiberServer *fiber.FiberServer
}

func NewWebServer(cfg *config.Config) *WebServer {
	logger.InfoOutCtx("initializing webserver SUCCESSFULLY")

	_, cancel := context.WithTimeout(context.Background(), cfg.Database.ConnTimeout)
	defer cancel()

	application := application.NewApplication()

	fiberServer := fiber.NewFiberServer(cfg)
	fiberServer.SetupRoutes(application)
	return &WebServer{
		cfg: cfg,
		fiberServer: fiberServer,
	}
}

func (s *WebServer) Run() {
	logger.InfoOutCtx("starting fiber server on port: " + s.cfg.HTTP.Port)
	if err := s.fiberServer.App.Listen(":" + s.cfg.HTTP.Port); err != nil {
		logger.FatalOutCtx("failed to start HTTP server", zap.Error(err))
	}
}

func (s *WebServer) Shutdown() {
	logger.InfoOutCtx("webserver is shutting down SUCCESSFULLY")
}
