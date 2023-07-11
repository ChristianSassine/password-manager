package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// TODO: move somewhere more appropriate
type Credentials struct {
	Username string
	Password string
}

type keysChange struct {
	OldKey string `json:"oldKey"`
	NewKey string `json:"newKey"`
}

const (
	authPath     = "/auth"
	passwordPath = "/password"
)

const (
	PassKeyQuery = "?key="
)

var client = &http.Client{}

func GetPassword(key string) (*http.Response, error) {
	url, err := getURL()
	if err != nil {
		return nil, err
	}
	response, err := client.Get(url + passwordPath + PassKeyQuery + key)
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
	response, err := client.Post(url+passwordPath, "text/plain", strings.NewReader(key))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ChangePasswordKey(oldKey string, newKey string) (*http.Response, error) {
	url, err := getURL()
	if err != nil {
		return nil, err
	}

	change := &keysChange{OldKey: oldKey, NewKey: newKey}
	data, err := json.Marshal(change)
	if err != nil {
		return nil, err
	}

	fullURL := url + passwordPath
	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func RemovePassword(key string) (*http.Response, error) {
	url, err := getURL()
	if err != nil {
		return nil, err
	}

	fullURL := url + passwordPath + PassKeyQuery + key
	req, err := http.NewRequest(http.MethodDelete, fullURL, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
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
