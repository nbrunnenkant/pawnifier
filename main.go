package main

import (
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
		panic(err)
	}

	providers := make([]MailProvider, 0)
	providers = append(providers, simplelogin.NewSimpleloginService())

	server.StartServer()
}
