package appinfohandlers

import (
	"github.com/WhoaWicked/board-game-shop/config"
	appinfousecases "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoUsecases"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/pkg/shopauth"
	"github.com/gofiber/fiber/v3"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
)

type IAppinfoHandler interface {
	GenerateApiKey(c fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfousecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfousecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c fiber.Ctx) error {
	apiKey, err := shopauth.NewShopAuth(shopauth.ApiKey, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}
