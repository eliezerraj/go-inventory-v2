package main

import (
	"os"
	"os/signal"
	"syscall"
	"context"

	"go.uber.org/zap"
	stdLog "log"

	"github.com/go-inventory-v2/cmd/webserver"
	"github.com/go-inventory-v2/application/config"
	"github.com/eliezerraj/go-core/v3/logger"
)

// Setup logging configuration and initialize logger
func setupLogging(cfg *config.Config) {
	logger.NewLogger(
		cfg.Log.Level,
		cfg.Log.Mode,
	).WithHook(func(ctx context.Context) []zap.Field {
		fields := []zap.Field{}
		return fields
	})
}

// getCmd retrieves the command type from the environment variable or uses the default value.
func getCmd(env string, val string) string {
	cmd := os.Getenv(env)
	if cmd == "" {
		cmd = val
	}
	return cmd
}

func main() {
	// Load environment configurations
	cfg, err := config.Load()
	if err != nil {
		stdLog.Fatalf("load environment configurations error: %+v", err)
		return
	}
	
	setupLogging(cfg)
	defer logger.Close()

	logger.InfoOutCtx("starting application", 
		zap.String("app_name", cfg.App.Name), 
		zap.String("version", cfg.App.Version))
	logger.InfoOutCtx("application configuration", zap.Any("config", cfg))

	// Setup signal handling for graceful shutdown
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	cmd := getCmd("COMMAND_TYPE", "webserver")

	switch cmd {
	case "worker":
		logger.InfoOutCtx("worker process NOT implemented")
	case "webserver":
		logger.InfoOutCtx("starting webserver process")
		
		webServer := webserver.NewWebServer(cfg)
		go webServer.Run()

		<-stopSignal
		webServer.Shutdown()
	}
}