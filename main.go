package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func readFile(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("localfile.txt")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

func writeFile(w http.ResponseWriter, r *http.Request) {
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("localfile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(string(responseData)); err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("localfile.txt")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

func clearFile(w http.ResponseWriter, r *http.Request) {
	err := os.Truncate("localfile.txt", 0)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile("localfile.txt")
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Get("/read", readFile)
	r.Post("/write", writeFile)
	r.Post("/clear", clearFile)
	http.ListenAndServe(":"+port, r)
}