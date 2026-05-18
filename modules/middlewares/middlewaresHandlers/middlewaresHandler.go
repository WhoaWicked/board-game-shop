package middlewareshandlers

import (
	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	middlewaresusecases "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresUsecases"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type middlewaresHandlersErrCode string

const (
	routerCheckErr middlewaresHandlersErrCode = "middleware-001"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	Logger() fiber.Handler
	RouterCheck() fiber.Handler
}

type middlewaresHandler struct {
	middlewaresUsecase middlewaresusecases.IMiddlewaresUsecase
	cfg                config.IConfig
}

func MiddlewaresHandler(middlewaresUsecase middlewaresusecases.IMiddlewaresUsecase, cfg config.IConfig) IMiddlewaresHandler {
	return &middlewaresHandler{
		middlewaresUsecase: middlewaresUsecase,
		cfg:                cfg,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{},
		AllowCredentials: false,
		ExposeHeaders:    []string{},
		MaxAge:           0,
	})
}
func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}
