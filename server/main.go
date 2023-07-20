package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ChristianSassine/password-manager/server/internal/handlers"
	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
)

var (
	ErrServerPort = errors.New("missing SERVER_PORT environmental variable for the server listening port")
)

func main() {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		log.Fatal(ErrServerPort)
	}

	// Start MongoDB Client
	log.Println("Starting Database...")
	err := mongodb.Start()
	if err != nil {
		log.Fatal("MongoDB: ", err)
	}
	defer mongodb.Disconnect()

	// Server initialization
	log.Println("Started Listening...")
	handlers.InitHandlers()

	serverPort := fmt.Sprintf(":%v", port)
	err = http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
