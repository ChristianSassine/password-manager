package requests

import (
	"fmt"
	"net/http"
	"strings"
)

// TODO: move somewhere more appropriate
type Credentials struct {
	Username string
	Password string
}

const (
	authPath     = "/auth"
	passwordPath = "/password"
)

const (
	PassKeyquery = "?key="
)

func GetPassword(key string) (*http.Response, error) {
	url, err := getURL()
	if err != nil {
		return nil, err
	}
	response, err := http.Get(url + passwordPath + PassKeyquery + key)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AddPassword(key string) (*http.Response, error) {
	url, err := getURL()
	if err != nil {
		return nil, err
	}
	response, err := http.Post(url+passwordPath, "text/plain", strings.NewReader(key))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getUserCreds() (Credentials, error) {
	// TODO: add logic
	return Credentials{"hello", "lol"}, nil
}

func getURL() (string, error) {
	creds, err := getUserCreds()
	if err != nil {
		return "", err
	}
	var link = fmt.Sprintf("http://%s:%s@localhost:4200", creds.Username, creds.Password) // TODO: load link from file

	return link, nil
}
