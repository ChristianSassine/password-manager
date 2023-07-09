package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChristianSassine/password-manager/server/internal/manager"
)

func handleAuth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var creds credentials
		err := decoder.Decode(&creds)
		if err != nil {
			panic(err)
		}
		err = manager.CreateUser(creds.Username, creds.Password)
		switch err {
		case nil:
			log.Println(creds, "Accepted") // TODO: remove
			w.WriteHeader(http.StatusAccepted)
		case manager.UserAlreadyExistsErr:
			log.Println(creds, "StatusConflict") // TODO: remove
			w.WriteHeader(http.StatusConflict)
		default:
			log.Println(creds, "StatusInternalServerError") // TODO: remove
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.Error(w, noMethodMsg, http.StatusForbidden)
	}
}
