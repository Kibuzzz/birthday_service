package main

import (
	"birtday_service/cron"
	"birtday_service/db"
	"birtday_service/notification"
	"birtday_service/repository/postgres"
	"birtday_service/server"
	"log"

	"go.uber.org/zap"
)

func main() {

	db, err := db.New()
	if err != nil {
		panic(err)
	}

	userStore := postgres.NewPostgresUsers(db)
	subStore := postgres.NewPostgresSubs(db)
	notificator := notification.New(userStore)
	scheduler := cron.New(subStore, notificator)
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	go func() {
		scheduler.Start()
	}()

	server := server.New(userStore, subStore, scheduler, sugar)
	server.InitRoutes()
	log.Fatal(server.Start())
}
