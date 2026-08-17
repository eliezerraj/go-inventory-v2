package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/json-iterator/go"
	"github.com/go-inventory-v2/cmd/webserver/framework/fiber/adapter"

	"github.com/go-inventory-v2/application/config"
	"github.com/go-inventory-v2/application/infrastructure/application"

	"github.com/eliezerraj/go-core/v3/logger"
)

type FiberServerConfig struct {
	fiber.Config
	Port string
}

func NewServerConfig(cfg *config.HTTP) FiberServerConfig {
	logger.InfoOutCtx("initializing fiber server config SUCCESSFULLY")

	return FiberServerConfig{
		Config: fiber.Config{
			JSONEncoder:           jsoniter.Marshal,
			JSONDecoder:           jsoniter.Unmarshal,
			DisableStartupMessage: cfg.DisableStartupMessage,
			ReadBufferSize:        cfg.ReadBufferSize,
			ReadTimeout:           cfg.ReadTimeout,
			WriteTimeout:          cfg.WriteTimeout,
			IdleTimeout:           cfg.IdleConnTimeout,
		},
		Port: cfg.Port,
	}
}

type httpAdapter struct {
	metadataAdp  	*adapter.MetadataAdapter
	applicationAdp   *adapter.ApplicationAdapter
}

func newAdapters(cfg *config.Config, application *application.Application) *httpAdapter {
	logger.InfoOutCtx("initializing fiber adapters SUCCESSFULLY")

	return &httpAdapter{
		metadataAdp:   	adapter.NewMetadataAdapter(cfg),
		applicationAdp:  adapter.NewApplicationAdapter(cfg, application),
	}
}

type FiberServer struct {
	cfg *config.Config
	App    *fiber.App
	fiberConfig FiberServerConfig
}

func NewFiberServer(cfg *config.Config) *FiberServer {
	logger.InfoOutCtx("initializing fiber server SUCCESSFULLY")

	// Create Fiber server configuration
	fiberConfig := NewServerConfig(&cfg.HTTP)
	app := fiber.New(fiberConfig.Config)

	return &FiberServer{
		cfg:    cfg,
		App:    app,
		fiberConfig: fiberConfig,
	}
}

func (s *FiberServer) SetupRoutes(application *application.Application) {
	logger.InfoOutCtx("setting up routes for fiber server SUCCESSFULLY")

	root := s.App.Group("/")

	// Create adapters for controllers						
	adapters := newAdapters(s.cfg, application)
	root.Get("/health", adapters.metadataAdp.HealthGet)

	appRoutes := root.Group("/v1")
	appRoutes.Get("/info", adapters.metadataAdp.InfoGet)
	appRoutes.Get("/echo-header", adapters.metadataAdp.HeadersGet)
	appRoutes.Get("/echo-context", adapters.metadataAdp.ContextGet)
	appRoutes.Get("/product/:sku", adapters.applicationAdp.ProductGet)
	appRoutes.Post("/product", adapters.applicationAdp.ProductAdd)
}