package main

import (
	"log"

	"github.com/jamascrorpJS/auth/internal/delivery"
	"github.com/jamascrorpJS/auth/internal/repository"
	"github.com/jamascrorpJS/auth/internal/server"
	"github.com/jamascrorpJS/auth/internal/service"
	"github.com/jamascrorpJS/auth/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	c, err := database.NewMongoDb()
	if err != nil {
		log.Fatal(err)
	}
	repository := repository.NewRepository(*c.Database("database"))
	service := service.NewService(repository)
	handler := delivery.NewHandler(service)
	handler.Start()
	server := server.NewServer(":8080", handler.Mux)
	server.StartServer()
}
