package main

import (
	"log"
	"net/http"
)

func main() {

	if err := DB.start(); err != nil {
		log.Println(err)
		return
	}

	//err := createTable()
	//println(err)

	defer DB.stop()
	Handler()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
