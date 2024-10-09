package main

import (
	"log"
	"music/internal/base"
	"music/internal/config"
	"music/internal/service"
	_ "music/docs"

	httpSwagger "github.com/swaggo/http-swagger"
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

	service := service.NewService(config,repository)
	defer func(){
		if err := service.Close(); err != nil{
			log.Fatalln(err)
		}
	}()

    // Подключение Swagger UI
    service.Router().PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    if err := service.Run(); err != nil {
        log.Fatalln(err)
    }
}