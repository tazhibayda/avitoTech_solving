package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"username"`
		Dbname   string `yaml:"dbname"`
		Password string `yaml:"password"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

var DB *sqlx.DB

func ConfigInit(cnf Config) {
	var err error
	DB, err = initDB(cnf)
	if err != nil {
		panic(err)
	}
}

func initDB(cnf Config) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cnf.Database.Host, cnf.Database.Port, cnf.Database.User, cnf.Database.Dbname, cnf.Database.Password, cnf.Database.Sslmode))

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	return db, nil
}
