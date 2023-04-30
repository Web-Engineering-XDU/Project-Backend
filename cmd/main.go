package main

import (
	"github.com/Web-Engineering-XDU/Project-Backend/app/controller"
	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	"github.com/Web-Engineering-XDU/Project-Backend/config"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var config config.Config
	err:= cleanenv.ReadConfig("../config/config.yml", &config)
	if err != nil {
		panic(err)
	}
	
	huggo := service.New(config, controller.SetController)
	huggo.Run(":8080")
}