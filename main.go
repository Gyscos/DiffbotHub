package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func readToken() string {
	body, err := ioutil.ReadFile("token")
	if err != nil {
		log.Println("Could not read file \"token\". Does this file exist?")
		log.Fatal(err)
	}

	return strings.Trim(string(body), " \n")
}

func main() {
	token := readToken()

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		api := r.FormValue("api")
		url := r.FormValue("url")
		log.Printf("%v (%v)", url, api)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("index.html")
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, struct {
			Token string
		}{
			token,
		})
	})
	log.Println("Listing now on port 8888")
	log.Println(http.ListenAndServe(":8888", nil))
}
