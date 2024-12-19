package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	UserName string
	Password string
	Name     string
	Port     string
	AppName  string
	Extras   string
}

func NewPostgresDB(params Config) (*sqlx.DB, error) {
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s application_name=%s %s",
		params.Host,
		params.UserName,
		params.Password,
		params.Name,
		params.Port,
		params.AppName,
		params.Extras,
	)

	db, err := sqlx.Connect("postgres", dbUrl)

	return db, err
}
