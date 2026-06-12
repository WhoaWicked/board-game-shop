package gameshandlers

import (
	"strings"

	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/modules/games"
	gamesusecases "github.com/WhoaWicked/board-game-shop/modules/games/gamesUsecases"
	"github.com/gofiber/fiber/v3"
)

type gameHandlerErrCode string

const (
	findOneGameErr gameHandlerErrCode = "games-001"
	findGamesErr   gameHandlerErrCode = "games-002"
	addGamesErr    gameHandlerErrCode = "games-003"
	deleteGamesErr gameHandlerErrCode = "games-004"
	updateGamesErr gameHandlerErrCode = "games-005"
)

type IGamesHandler interface {
	FindOneGame(c fiber.Ctx) error
	FindGames(c fiber.Ctx) error
	AddGame(c fiber.Ctx) error
	DeleteGame(c fiber.Ctx) error
	UpdateGame(c fiber.Ctx) error
}

type gamesHandler struct {
	gamesUsecase gamesusecases.IGamesUsecase
}

func GamesHandler(gamesUsecase gamesusecases.IGamesUsecase) IGamesHandler {
	return &gamesHandler{
		gamesUsecase: gamesUsecase,
	}
}

func (h *gamesHandler) FindOneGame(c fiber.Ctx) error {
	gameId := strings.Trim(c.Params("game_id"), " ")
	game, err := h.gamesUsecase.FindOneGame(gameId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneGameErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, game).Res()
}

func (h *gamesHandler) FindGames(c fiber.Ctx) error {
	req := &games.GameFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
		Categories:    make([]int, 0),
	}
	if err := c.Bind().Query(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findGamesErr),
			err.Error(),
		).Res()
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}
	cleanCategories := make([]int, 0)
	for _, catId := range req.Categories {
		if catId > 0 {
			cleanCategories = append(cleanCategories, catId)
		}
	}
	req.Categories = cleanCategories
	games := h.gamesUsecase.FindGames(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, games).Res()
}

func (h *gamesHandler) AddGame(c fiber.Ctx) error {
	req := &games.Game{
		Category: make([]*appinfo.Category, 0),
		Image:    make([]*entities.Images, 0),
	}
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addGamesErr),
			err.Error(),
		).Res()
	}
	game, err := h.gamesUsecase.AddGame(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addGamesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, game).Res()
}

func (h *gamesHandler) DeleteGame(c fiber.Ctx) error {
	gameId := strings.Trim(c.Params("game_id"), " ")
	if err := h.gamesUsecase.DeleteGame(gameId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteGamesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, "delete success").Res()
}

func (h *gamesHandler) UpdateGame(c fiber.Ctx) error {
	gameId := strings.Trim(c.Params("game_id"), " ")
	req := &games.Game{
		Category: make([]*appinfo.Category, 0),
		Image:    make([]*entities.Images, 0),
	}
	if err := c.Bind().Body(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateGamesErr),
			err.Error(),
		).Res()
	}
	req.Id = gameId
	game, err := h.gamesUsecase.UpdateGame(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateGamesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, game).Res()
}
