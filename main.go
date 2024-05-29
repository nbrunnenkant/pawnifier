package main

import (
	"github.com/joho/godotenv"
	"github.com/nbrunnenkant/pawnifier/hibp"
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

	emails := make([]string, 0)
	for _, provider := range providers {
		emails = append(emails, provider.GetMails()...)
	}

	hibp.CheckMails(emails...)
}
