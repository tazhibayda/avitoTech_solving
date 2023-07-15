package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tazhibayda/avitoTech_solving/pkg/database"
	"github.com/tazhibayda/avitoTech_solving/pkg/handler"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func main() {

	cnf, err := os.ReadFile("configs/config.yml")

	fmt.Println(cnf)
	if err != nil {
		panic(err)
	}

	var config database.Config
	var api struct {
		ExchangeApi string `yaml:"exchange-api"`
	}
	err = yaml.Unmarshal(cnf, &config)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(cnf, &api)

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
