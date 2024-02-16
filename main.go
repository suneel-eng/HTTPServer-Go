package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func logger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		f(w, r)
	}
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func main() {

	r := mux.NewRouter()

	// r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Fprintf(w, "Welcome! You are at %s", req.URL.Path)
	// })

	// r.HandleFunc("/number/{num}", func(w http.ResponseWriter, req *http.Request) {
	// 	vars := mux.Vars(req)
	// 	fmt.Fprintf(w, "Your number is %s", vars["num"])
	// })

	// fs := http.FileServer(http.Dir("./static/public/"))

	// r.Handle("/", http.StripPrefix("/", fs))

	r.HandleFunc("/", logger(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "cookie-username")

		if r.Method == http.MethodGet {

			if value, ok := session.Values["username"]; !ok {
				fmt.Println(ok)
				tmpl := template.Must(template.ParseFiles("./static/public/index.html"))
				tmpl.Execute(w, nil)
			} else {
				fmt.Println(ok)
				username := value
				tmpl := template.Must(template.ParseFiles("./static/private/dashboard.html"))
				tmpl.Execute(w, username)
			}

			return
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			session.Values["username"] = username
			session.Save(r, w)
			tmpl := template.Must(template.ParseFiles("./static/private/dashboard.html"))
			tmpl.Execute(w, username)
		}
	})).Methods("POST", "GET")

	r.HandleFunc("/logout", logger(func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "cookie-username")
		session.Options.MaxAge = -1
		session.Save(req, w)
	})).Methods("POST")

	http.ListenAndServe("localhost:3000", r)
}
