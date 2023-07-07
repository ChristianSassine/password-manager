package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChristianSassine/password-manager/server/internal/handlers"
	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()

	// Start MongoDB Client
	log.Println("Starting Database...")
	err := mongodb.Start()
	if err != nil {
		log.Fatal("MongoDB: ", err) // TODO: replace log
	}
	defer mongodb.Disconnect()

	// Server initialization
	log.Println("Started Listening...")
	handlers.InitHandlers()
	err = http.ListenAndServe(":4200", nil) // TODO: Change port to environment PORT
	if err != nil {
		log.Fatal("ListenAndServe: ", err) // TODO: replace log
	}
}
