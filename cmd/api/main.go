package main

import (
	"log"

	"github.com/mykhalskyio/image-api/internal/config"
	"github.com/mykhalskyio/image-api/internal/controller/http"
	"github.com/mykhalskyio/image-api/internal/service"
	"github.com/mykhalskyio/image-api/internal/storage/postgres"
)

func main() {
	config := config.GetConfig()
	db, err := postgres.NewPostgres(config)
	if err != nil {
		log.Fatalln("Connect to db :", err)
	}
	imageService := service.NewImageService(db, config)
	router := http.NewRouter(imageService)
	router.Run(":8080")
}
