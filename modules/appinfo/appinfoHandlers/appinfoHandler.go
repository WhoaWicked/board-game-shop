package appinfohandlers

import (
	"github.com/WhoaWicked/board-game-shop/config"
	appinfousecases "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoUsecases"
)

type IAppinfoHandler interface{}

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
