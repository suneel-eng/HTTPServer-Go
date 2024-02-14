package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func logger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		f(w, r)
	}
}

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

		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("./static/public/index.html"))
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			tmpl := template.Must(template.ParseFiles("./static/private/dashboard.html"))
			tmpl.Execute(w, username)
		}
	})).Methods("POST", "GET")

	http.ListenAndServe("localhost:3000", r)
}
