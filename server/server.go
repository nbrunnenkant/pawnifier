package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/nbrunnenkant/pawnifier/hibp"
	"github.com/nbrunnenkant/pawnifier/simplelogin"
)

type Mails struct {
	Mails []string
}

var hibpService hibp.HIBPService

func StartServer() {
	hibpService = *hibp.NewHIBPService()

	fsys := os.DirFS("server/views/static")
	fs := http.FileServerFS(fsys)

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("GET /status/{mail}", handleMailStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("server/views/index.html")

	if err != nil {
		fmt.Println(err)
	}

	sl := simplelogin.NewSimpleloginService()
	checkingMails := sl.GetMails()
	mails := Mails{Mails: checkingMails}
	err = tmpl.Execute(w, mails)

	if err != nil {
		fmt.Println(err)
	}
}

func handleMailStatus(w http.ResponseWriter, r *http.Request) {
	hibpService.AddMail(r)
	test := <-hibpService.Response

	var cool string
	if test {
		cool = "jau is sicher"
	} else {
		cool = "nicht so sicher"
	}
	w.Write([]byte(cool))
}
