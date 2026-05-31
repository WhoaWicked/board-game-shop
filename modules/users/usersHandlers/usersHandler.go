package usershandlers

import (
	"strings"

	"github.com/WhoaWicked/board-game-shop/config"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/modules/users"
	usersusecases "github.com/WhoaWicked/board-game-shop/modules/users/usersUsecases"
	"github.com/WhoaWicked/board-game-shop/pkg/shopauth"
	"github.com/gofiber/fiber/v3"
)

type userHandlersErrCode string

const (
	signUpCustomerErr  userHandlersErrCode = "user-001"
	signUpAdminErr     userHandlersErrCode = "user-002"
	signInErr          userHandlersErrCode = "user-003"
	refreshPassportErr userHandlersErrCode = "user-004"
	deleteOauthErr     userHandlersErrCode = "user-005"
	generateAdminToken userHandlersErrCode = "user-006"
	getUserProfileErr  userHandlersErrCode = "user-007"
)

type IUsersHandler interface {
	InsertCustomer(c fiber.Ctx) error
	InsertAdmin(c fiber.Ctx) error
	SignIn(c fiber.Ctx) error
	RefreshPassport(c fiber.Ctx) error
	SignOut(c fiber.Ctx) error
	GenerateAdminToken(c fiber.Ctx) error
	GetUserProfile(c fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersusecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersusecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) InsertCustomer(c fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			err.Error(),
		).Res()
	}
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			"email is invalid",
		).Res()
	}
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
func (h *usersHandler) InsertAdmin(c fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			err.Error(),
		).Res()
	}
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			"email is invalid",
		).Res()
	}
	result, err := h.usersUsecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteOauthErr),
			err.Error(),
		).Res()
	}
	if err := h.usersUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteOauthErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GenerateAdminToken(c fiber.Ctx) error {
	adminToken, err := shopauth.NewShopAuth(shopauth.Admin, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(generateAdminToken),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *usersHandler) GetUserProfile(c fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")
	profile, err := h.usersUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, profile).Res()
}
