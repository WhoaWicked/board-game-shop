package server

import (
	middlewareshandlers "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/WhoaWicked/board-game-shop/modules/middlewares/middlewaresUsecases"
	monitorhandlers "github.com/WhoaWicked/board-game-shop/modules/monitor/monitorHandlers"
	usershandlers "github.com/WhoaWicked/board-game-shop/modules/users/usersHandlers"
	usersrepositories "github.com/WhoaWicked/board-game-shop/modules/users/usersRepositories"
	usersusecases "github.com/WhoaWicked/board-game-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v3"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
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

func (m *moduleFactory) UsersModule() {
	repository := usersrepositories.UsersRepository(m.s.db)
	usecase := usersusecases.UsersUsecase(m.s.cfg, repository)
	handler := usershandlers.UsersHandler(m.s.cfg, usecase)
	router := m.r.Group("/users")
	router.Post("/signup", handler.InsertCustomer)
	router.Post("/signup-admin", handler.InsertAdmin)
	router.Post("/signin", handler.SignIn)
}
