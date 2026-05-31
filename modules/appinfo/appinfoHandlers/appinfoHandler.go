package appinfohandlers

import (
	"strconv"
	"strings"

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
	addCategoryErr    appinfoHandlersErrCode = "appinfo-003"
	removeCategoryErr appinfoHandlersErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c fiber.Ctx) error
	FindCategories(c fiber.Ctx) error
	AddCategory(c fiber.Ctx) error
	RemoveCategory(c fiber.Ctx) error
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

func (h *appinfoHandler) AddCategory(c fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	if err := h.appinfoUsecase.AddCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) RemoveCategory(c fiber.Ctx) error {
	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"id type is invalid",
		).Res()
	}
	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"id must more than 0",
		).Res()
	}
	if err := h.appinfoUsecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
