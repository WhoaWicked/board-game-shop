package server

import (
	appinfohandlers "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoHandlers"
	appinforepositories "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoRepositories"
	appinfousecases "github.com/WhoaWicked/board-game-shop/modules/appinfo/appinfoUsecases"
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
	AppinfoModule()
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
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(4), handler.GenerateAdminToken)
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinforepositories.AppinfoRepository(m.s.db)
	usecase := appinfousecases.AppinfoUsecase(repository)
	handler := appinfohandlers.AppinfoHandler(m.s.cfg, usecase)
	router := m.r.Group("/appinfo")
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(4), handler.GenerateApiKey)
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategories)
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(4), handler.AddCategory)
	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(4), handler.RemoveCategory)
}
