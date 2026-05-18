package server

import (
	middlewareshandlers "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresUsecases"
	monitorhandlers "github.com/WhoaWicked/board-game-shop/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewareshandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewareshandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewareshandlers.IMiddlewaresHandler {
	repository := middlewaresrepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresusecases.MiddlewaresUsecase(repository)
	handler := middlewareshandlers.MiddlewaresHandler(usecase, s.cfg)
	return handler
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorhandlers.MoniterHandler(m.s.cfg)
	m.r.Get("/", handler.HealthCheck)
}
