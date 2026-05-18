package monitorhandlers

import (
	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/modules/monitor"
	"github.com/gofiber/fiber/v3"
)

type IMonitorHandler interface {
	HealthCheck(c fiber.Ctx) error
}

type monitorHandler struct {
	cfg config.IConfig
}

func MoniterHandler(cfg config.IConfig) IMonitorHandler {
	return &monitorHandler{
		cfg: cfg,
	}
}

func (m *monitorHandler) HealthCheck(c fiber.Ctx) error {
	res := &monitor.Moniter{
		Name:    m.cfg.App().Name(),
		Version: m.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}
