package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DbHandler struct {
	db *sql.DB
}

func (handler *DbHandler) InsertApiKey(apiKey string) error {
	_, err := handler.db.Exec("insert into api_keys (api_key) values (?)", apiKey)
	return err
}

func NewHandler() *DbHandler {
	cfg := mysql.Config{
		User:                 "test",
		Passwd:               "example",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "pawnifier",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Error while processing database config", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("Could not connect to database", pingErr)
	}

	err = handleDatabaseMigration(db)
	if err != nil {
		log.Fatal("Could not process database migration", err)
	}

	return &DbHandler{db: db}
}

func handleDatabaseMigration(db *sql.DB) error {
	_, err := db.Exec("create table if not exists api_keys (id int auto_increment not null, api_key text, constraint pk_api_key primary key (id));")
	return err
}
