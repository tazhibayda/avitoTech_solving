package main

import (
	"avitoTech_solving/pkg/repository"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

func main() {

	cnf, err := os.ReadFile("configs/config.yml")

	if err != nil {
		panic(err)
	}

	var config repository.Config

	err = yaml.Unmarshal(cnf, &config)

	if err != nil {
		panic(err)
	}

	db, err := repository.InitDB(config)

	if err != nil {
		panic(err)
	}

	fmt.Println(db)
}
