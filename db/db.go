package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type DbHandler struct {
	db *sql.DB
}

func (handler *DbHandler) InsertApiKey(apiKey string) error {
	_, err := handler.db.Exec("insert into api_keys (api_key) values (?)", apiKey)
	return err
}

func (handler *DbHandler) GetApiKeys() []string {
	rows, err := handler.db.Query("select api_key from api_keys")
	apiKeys := make([]string, 0)

	if err != nil {
		log.Println("Could not retrieve api-keys", err)
		return apiKeys
	}

	for {
		if !rows.Next() {
			break
		}

		var apiKey string
		rows.Scan(&apiKey)
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys
}

func NewHandler() *DbHandler {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PW"),
		DBName:               os.Getenv("DB_NAME"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
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
