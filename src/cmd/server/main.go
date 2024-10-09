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
		log.Fatalln(err)
	}

	repository, err := base.NewRepository(config)
	if err != nil{
		log.Fatalln(err)
	}

	service := service.NewService(repository)
	defer func(){
		if err := service.Close(); err != nil{
			log.Fatalln(err)
		}
	}()

	if err := service.Run(); err != nil{
		log.Fatalln(err)
	}
}