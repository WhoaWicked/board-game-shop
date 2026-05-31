package middlewareshandlers

import (
	"strings"

	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	middlewaresusecases "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresUsecases"
	"github.com/WhoaWicked/board-game-shop/pkg/shopauth"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type middlewaresHandlersErrCode string

const (
	routerCheckErr middlewaresHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewaresHandlersErrCode = "middleware-002"
	paramsCheckErr middlewaresHandlersErrCode = "middleware-003"
	authorizeErr   middlewaresHandlersErrCode = "middleware-004"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	Logger() fiber.Handler
	RouterCheck() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRole ...int) fiber.Handler
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

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		// ดึง token
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := shopauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}
		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.Role_id)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 4 {
			c.Next()
		}
		if userId != c.Params("user_id") {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRole ...int) fiber.Handler {
	return func(c fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"role id is not int type",
			).Res()
		}
		combinedRole := 0
		for _, role := range expectRole {
			combinedRole |= role
		}
		if (userRoleId & combinedRole) != 0 {
			return c.Next()
		}
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}
