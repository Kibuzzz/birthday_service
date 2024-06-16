package server

import (
	"birtday_service/cron"
	subs "birtday_service/models/subscription"
	user "birtday_service/models/user"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	UserStorage user.UserStorage
	SubsStorage subs.SubStorage
	Scheduler   cron.Cron
	App         *fiber.App
}

func New(users user.UserStorage, subs subs.SubStorage, schehuler cron.Cron) Server {
	app := fiber.New()
	return Server{UserStorage: users, SubsStorage: subs, Scheduler: schehuler, App: app}
}

func (s *Server) InitRoutes() {
	app := fiber.New()
	api := app.Group("/api")
	// employees
	api.Use(s.AuthMiddleware)
	api.Get("/employees", s.AllEmployees)
	api.Delete("/employees/:email", s.DeleteEmployee)
	// subs
	api.Get("/subs", s.SubsByID)       // все подписки пользователя
	api.Post("/subs", s.Subscribe)     // создание подписки
	api.Delete("/subs", s.UnSubscribe) // удаление подписки
	// authorization
	app.Post("/register", s.Register)
	app.Post("/login", s.Login)
	app.Post("/logout", s.Logout)
	s.App = app
}

func (s *Server) Start() error {
	return s.App.Listen(":1234")
}
