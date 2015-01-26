package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Frees up filename, by moving it into filebase.level
func logShift(filename string, filebase string, level int) error {
	newPath := fmt.Sprintf("%v.%v", filebase, level)
	if exists(newPath) {
		err := logShift(newPath, filebase, level+1)
		if err != nil {
			return err
		}
	}
	return os.Rename(filename, newPath)
}

func logRotate(history []string, filename string) error {
	if exists(filename) {
		err := logShift(filename, filename, 0)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(filename + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()

	for s := range history {
		fmt.Fprintln(f, s)
	}

	return nil
}

func main() {
	token := readToken()

	// history := make([]string, 0)
	historyChannel := make(chan string, 1)
	go func() {
		for s := range historyChannel {
			// history = append(history, s)
			log.Println(s)
		}
	}()

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		url := r.FormValue("url")
		historyChannel <- url
		fmt.Fprint(w, "OK")
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
	close(historyChannel)
}
