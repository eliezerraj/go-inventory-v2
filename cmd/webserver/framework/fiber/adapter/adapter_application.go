package adapter

import (
	"github.com/go-inventory-v2/application/config"
	//"github.com/go-inventory-v2/application/infrastructure/application"
	"github.com/eliezerraj/go-core/v3/logger"
)

type ApplicationAdapter struct {
	cfg *config.Config
}

func NewApplicationAdapter(cfg *config.Config) *ApplicationAdapter {
	logger.InfoOutCtx("initializing application adapter SUCCESSFULLY")

	return &ApplicationAdapter{
		cfg: cfg,
	}
}
