package dataBase

import (
	"Cloud/logger"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type connectDB struct {
	DB *sql.DB
}

func ConnectDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=1234 dbname=postgres sslmode=disable" // возможно сделать функцию которая эти параметры принимает будет лучшей идеей!
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Failed connect to database!" + err.Error())
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logger.Error("Failed ping to database!" + err.Error())
		log.Fatal(err)
	}
	return db
}
