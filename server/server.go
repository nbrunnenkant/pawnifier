package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/nbrunnenkant/pawnifier/simplelogin"
)

type test struct {
	Mails []string
}

func StartServer() {
	fsys := os.DirFS("server/views/static")
	fs := http.FileServerFS(fsys)

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("server/views/index.html")

	if err != nil {
		fmt.Println(err)
	}

	test := test{Mails: simplelogin.NewSimpleloginService().GetMails()}
	_ = tmpl.Execute(w, test)
}
