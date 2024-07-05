package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nbrunnenkant/pawnifier/hibp"
)

type Mails struct {
	Mails []string
}

var hibpService hibp.HIBPService

func StartServer() {
	hibpService = *hibp.NewHIBPService()

	fsys := os.DirFS("server/views/static")
	fs := http.FileServerFS(fsys)

	http.Handle("GET /static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("GET /", handleIndex)
	http.HandleFunc("GET /login", handleLogin)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("server/views/index.html")

	if err != nil {
		log.Print(err)
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		fmt.Println(err)
	}
}

type User struct {
	Email string `json:"email"`
}

type Response struct {
	User User `json:"user"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tokenUrl := url.URL{
		Scheme: "https",
		Host:   "app.simplelogin.io",
		Path:   "oauth2/token",
	}

	query := url.Values{}
	query.Set("grant_type", "authorization_code")
	query.Set("code", code)
	query.Set("redirect_url", "http://localhost:8080/login")
	query.Set("client_id", "pawnifier-rvvtvfkzyf")
	query.Set("client_secret", os.Getenv("SIWSL_CLIENT_SECRET"))

	req, err := http.NewRequest(http.MethodPost, tokenUrl.String(), strings.NewReader(query.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, _ := client.Do(req)

	bytes, _ := io.ReadAll(res.Body)
	fmt.Println(string(bytes))

	user := &Response{}
	json.Unmarshal(bytes, user)

	fmt.Println(user)
	tmpl, err := template.ParseFiles("server/views/index.html")

	if err != nil {
		log.Print(err)
	}

	err = tmpl.Execute(w, user.User)

	if err != nil {
		fmt.Println(err)
	}
}
