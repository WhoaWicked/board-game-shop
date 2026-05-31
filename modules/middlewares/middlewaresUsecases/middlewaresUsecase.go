package middlewaresusecases

import middlewaresrepositories "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
	FindAccessToken(userId string, accessToken string) bool
}

type middlewaresUsecase struct {
	middlewaresRepository middlewaresrepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewaresRepository middlewaresrepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}
