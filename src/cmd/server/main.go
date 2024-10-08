package main

import (
	"log"
	"music/internal/base"
	"music/internal/config"
	"music/internal/service"
)

func main(){
	config, err := config.NewConfig()
	if err != nil{
		log.Fatalf("error loading .env", err)
	}

	repository, err := base.NewRepository(config)
	if err != nil{
		log.Fatal(err)
	}

	service := service.NewService(repository)

	service.Run()
}