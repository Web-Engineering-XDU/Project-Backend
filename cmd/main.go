package main

import (
	"github.com/Web-Engineering-XDU/Project-Backend/app/controller"
	"github.com/Web-Engineering-XDU/Project-Backend/app/service"
	"github.com/Web-Engineering-XDU/Project-Backend/config"
	"github.com/ilyakaznacheev/cleanenv"

	_ "github.com/Web-Engineering-XDU/Project-Backend/docs/swaggo"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var config config.Config
	err:= cleanenv.ReadConfig("../config/config.yml", &config)
	if err != nil {
		panic(err)
	}
	
	huggo := service.New(config, controller.SetController)
	huggo.Run(":8080")
}