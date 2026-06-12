package gamesusecases

import (
	"math"

	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/modules/games"
	gamesrepositories "github.com/WhoaWicked/board-game-shop/modules/games/gamesRepositories"
)

type IGamesUsecase interface {
	FindOneGame(gameId string) (*games.Game, error)
	FindGames(req *games.GameFilter) *entities.PagianateRes
	AddGame(req *games.Game) (*games.Game, error)
	DeleteGame(gameId string) error
	UpdateGame(req *games.Game) (*games.Game, error)
}

type gamesUsecase struct {
	gamesRepository gamesrepositories.IGamesRepository
}

func GamesUsecase(gamesRepository gamesrepositories.IGamesRepository) IGamesUsecase {
	return &gamesUsecase{
		gamesRepository: gamesRepository,
	}
}

func (u *gamesUsecase) FindOneGame(gameId string) (*games.Game, error) {
	game, err := u.gamesRepository.FindOneGame(gameId)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (u *gamesUsecase) FindGames(req *games.GameFilter) *entities.PagianateRes {
	games, count := u.gamesRepository.FindGames(req)
	return &entities.PagianateRes{
		Data:      games,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *gamesUsecase) AddGame(req *games.Game) (*games.Game, error) {
	game, err := u.gamesRepository.InsertGame(req)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (u *gamesUsecase) DeleteGame(gameId string) error {
	if err := u.gamesRepository.DeleteGame(gameId); err != nil {
		return err
	}
	return nil
}

func (u *gamesUsecase) UpdateGame(req *games.Game) (*games.Game, error) {
	game, err := u.gamesRepository.UpdateGame(req)
	if err != nil {
		return nil, err
	}
	return game, nil
}
