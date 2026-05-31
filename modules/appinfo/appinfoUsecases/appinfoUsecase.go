package appinfousecases

import (
	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	appinforepositories "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	AddCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
	UpdateCategory(req *appinfo.Category) error
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

func (u *appinfoUsecase) AddCategory(req []*appinfo.Category) error {
	if err := u.appinfoRepository.InsertCategory(req); err != nil {
		return err
	}
	return nil
}

func (u *appinfoUsecase) UpdateCategory(req *appinfo.Category) error {
	if err := u.appinfoRepository.UpdateCategory(req); err != nil {
		return err
	}
	return nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.appinfoRepository.DeleteCategory(categoryId); err != nil {
		return err
	}
	return nil
}
