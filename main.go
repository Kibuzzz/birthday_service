package main

import (
	"birtday_service/cron"
	"birtday_service/db"
	"birtday_service/notification"
	"birtday_service/repository/postgres"
	"birtday_service/server"
	"log"
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

	go func() {
		scheduler.Start()
	}()

	server := server.New(userStore, subStore, scheduler)
	server.InitRoutes()
	log.Fatal(server.Start())
}
