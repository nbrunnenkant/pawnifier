package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nbrunnenkant/pawnifier/server"
	"github.com/nbrunnenkant/pawnifier/simplelogin"
)

type MailProvider interface {
	GetMails() []string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	providers := make([]MailProvider, 0)
	providers = append(providers, simplelogin.NewSimpleloginService())

	server.StartServer()
}
