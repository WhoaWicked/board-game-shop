package appinfohandlers

import (
	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	appinfousecases "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoUsecases"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/pkg/shopauth"
	"github.com/gofiber/fiber/v3"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
	findCategoriesErr appinfoHandlersErrCode = "appinfo-002"
)

type IAppinfoHandler interface {
	GenerateApiKey(c fiber.Ctx) error
	FindCategories(c fiber.Ctx) error
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

func (h *appinfoHandler) FindCategories(c fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoriesErr),
			err.Error(),
		).Res()
	}
	categories, err := h.appinfoUsecase.FindCategories(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoriesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, categories).Res()
}
