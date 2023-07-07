package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ChristianSassine/password-manager/server/internal/manager"
)

var (
	WrongAuthTypeErr = errors.New("Invalid authentification. The authorization header should use the Basic authentification")
	InvalidAuthErr   = errors.New("Invalid authentification. The authorization should follow the following pattern <username>:<password>")
)

const noMethodMsg = "Method not allowed"
const credsPattern = "^[^:\n]+:[^:\n]+$"

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func InitHandlers() {
	http.HandleFunc("/auth", postAccount)
}

func postAccount(w http.ResponseWriter, r *http.Request) {
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
			log.Println(creds, "Accepted")
			w.WriteHeader(http.StatusAccepted)
		case manager.UserAlreadyExistsErr:
			log.Println(creds, "StatusConflict")
			w.WriteHeader(http.StatusConflict)
		default:
			log.Println(creds, "StatusInternalServerError")
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		http.Error(w, noMethodMsg, http.StatusMethodNotAllowed)
	}
}

func getCredentials(r *http.Request) (credentials, error) {
	auth := r.Header.Get("Authorization")
	correctType := strings.HasPrefix(auth, "Basic ")
	if !correctType {
		return credentials{}, WrongAuthTypeErr
	}
	auth = strings.TrimPrefix(auth, "Basic ")
	credsBytes, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return credentials{}, err
	}
	decodedCreds := string(credsBytes)
	match, err := regexp.MatchString(credsPattern, decodedCreds)
	if err != nil {
		return credentials{}, err
	}
	if !match {
		return credentials{}, InvalidAuthErr
	}

	creds := strings.Split(decodedCreds, ":")
	username, password := creds[0], creds[1]

	return credentials{Username: username, Password: password}, nil
}
