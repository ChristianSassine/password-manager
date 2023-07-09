package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ChristianSassine/password-manager/server/internal/manager"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func handlePassword(w http.ResponseWriter, r *http.Request) {
	creds, err := getCredentials(r)
	if err != nil {
		log.Fatal("Creds: ", err) // TODO: Return an error instead of panicking
	}
	log.Println("Creds Accepted", creds) // TODO: remove
	switch r.Method {
	case http.MethodPost:
		handlePasswordAdd(w, r, creds)
	case http.MethodGet:
		handlePasswordGet(w, r, creds)
	case http.MethodDelete:
		handlePasswordRemove(w, r, creds)
	case http.MethodPut:
		handlePasswordModify(w, r, creds)
	default:
		http.Error(w, noMethodMsg, http.StatusMethodNotAllowed)
	}
}

func handlePasswordAdd(w http.ResponseWriter, r *http.Request, creds credentials) {

	log.Println("handlePasswordAdd", "Adding...") // TODO: remove
	if r.Header.Get("Content-type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("handlePasswordAdd", "WRONG HEADER") // TODO: remove
		return
	}

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("handlePasswordAdd", "err:", err) // TODO: remove
		return
	}
	key := string(responseData)
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("handlePasswordAdd", "error:", err) // TODO: remove
		return
	}
	err = manager.UserAddPassword(creds.Username, creds.Password, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("handlePasswordAdd", "Error:", err) // TODO: remove
		return
	}
	log.Println("handlePasswordAdd", "Password added!") // TODO: remove
	w.WriteHeader(http.StatusCreated)
}

func handlePasswordGet(w http.ResponseWriter, r *http.Request, creds credentials) {
	passKey := r.URL.Query().Get("key")
	if passKey == "" {
		log.Println("handlePasswordGet", "Password Empty") // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("GOT KEY", passKey)
	pass, err := manager.UserGetPassword(creds.Username, creds.Password, passKey)
	if err != nil {
		log.Println("handlePasswordGet", "err", err) // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pass))
}

func handlePasswordRemove(w http.ResponseWriter, r *http.Request, creds credentials) {

}

func handlePasswordModify(w http.ResponseWriter, r *http.Request, creds credentials) {

}
