package main

import (
	"avitoTech_solving/pkg/database"
	"avitoTech_solving/pkg/handler"
	"github.com/gin-gonic/gin"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func main() {

	cnf, err := os.ReadFile("configs/config.yml")

	if err != nil {
		panic(err)
	}

	var config database.Config

	err = yaml.Unmarshal(cnf, &config)

	if err != nil {
		panic(err)
	}

	database.ConfigInit(config)

	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	handler.Routers(router)
	router.Run()
}
