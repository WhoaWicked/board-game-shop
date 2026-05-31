package appinfousecases

import appinforepositories "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoRepositories"

type IAppinfoUsecase interface{}

type appinfoUsecase struct {
	appinfoRepository appinforepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}
