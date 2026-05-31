package appinfousecases

import (
	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	appinforepositories "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoUsecase struct {
	appinfoRepository appinforepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

func (u *appinfoUsecase) FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	categories, err := u.appinfoRepository.FindCategories(req)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
