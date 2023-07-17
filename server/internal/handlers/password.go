package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ChristianSassine/password-manager/server/internal/manager"
)

type renameKeys struct {
	OldKey string `json:"oldKey"`
	NewKey string `json:"newKey"`
}

func handlePassword(w http.ResponseWriter, r *http.Request) {
	creds, err := getCredentials(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := manager.ValidateUserCreds(creds.Username, creds.Password); err != nil {
		log.Println("Creds refused", creds) // TODO: remove
		w.WriteHeader(http.StatusForbidden)
		return
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
	if r.Header.Get("Content-type") != "application/json" {
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
	var passOpts manager.PasswordOptions
	err = json.Unmarshal(responseData, &passOpts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("handlePasswordAdd", "error:", err) // TODO: remove
		return
	}

	if passOpts.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("handlePasswordAdd", "error:", err) // TODO: remove
		return
	}

	err = manager.UserAddPassword(creds.Username, passOpts)
	if err == manager.PasswordConflictErr {
		log.Println("handlePasswordAdd", "err", err) // TODO: remove
		w.WriteHeader(http.StatusConflict)
		return
	}
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
	fmt.Println("GOT KEY", passKey) // TODO: remove
	pass, err := manager.UserGetPassword(creds.Username, passKey)
	if err == manager.NoPasswordErr {
		log.Println("handlePasswordGet", "err", err) // TODO: remove
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("handlePasswordGet", "err", err) // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pass))
}

func handlePasswordRemove(w http.ResponseWriter, r *http.Request, creds credentials) {
	passKey := r.URL.Query().Get("key")
	if passKey == "" {
		log.Println("handlePasswordGet", "Password Empty") // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("GOT KEY", passKey) // TODO: remove
	err := manager.UserRemovePassword(creds.Username, passKey)
	if err != nil {
		log.Println("handlePasswordGet", "err", err) // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handlePasswordModify(w http.ResponseWriter, r *http.Request, creds credentials) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("handlePasswordModify", "bad type", r.Header.Get("Content-Type")) // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var keys renameKeys
	err := json.NewDecoder(r.Body).Decode(&keys)
	if err != nil {
		log.Println("handlePasswordModify", "bad json", keys) // TODO: remove
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(keys) // TODO: remove
	err = manager.UserRenamePassword(creds.Username, keys.OldKey, keys.NewKey)
	if err == manager.NoPasswordErr {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == manager.PasswordConflictErr {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
