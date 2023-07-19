package handlers

import (
	"encoding/json"
	"io/ioutil"
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
		w.WriteHeader(http.StatusForbidden)
		return
	}
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

	if r.Header.Get("Content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var passOpts manager.PasswordOptions
	err = json.Unmarshal(responseData, &passOpts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if passOpts.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = manager.UserAddPassword(creds.Username, passOpts)
	if err == manager.PasswordConflictErr {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func handlePasswordGet(w http.ResponseWriter, r *http.Request, creds credentials) {
	passKey := r.URL.Query().Get("key")
	if passKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pass, err := manager.UserGetPassword(creds.Username, passKey)
	if err == manager.NoPasswordErr {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pass))
}

func handlePasswordRemove(w http.ResponseWriter, r *http.Request, creds credentials) {
	passKey := r.URL.Query().Get("key")
	if passKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := manager.UserRemovePassword(creds.Username, passKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handlePasswordModify(w http.ResponseWriter, r *http.Request, creds credentials) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var keys renameKeys
	err := json.NewDecoder(r.Body).Decode(&keys)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
