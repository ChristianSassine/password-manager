package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
)

var (
	NoUsernameErr = errors.New("unable to find the username environmental variable")
	NoPasswordErr = errors.New("unable to find the password environmental variable")
	NoURL         = errors.New("unable to locate the URL. Please use config to define the URL")
)

type keysChange struct {
	OldKey string `json:"oldKey"`
	NewKey string `json:"newKey"`
}

const (
	passwordPath = "/password"
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
	username, ok := os.LookupEnv("PASS_USERNAME")
	if !ok {
		output.Error(NoUsernameErr.Error())
		os.Exit(1)
	}

	password, ok := os.LookupEnv("PASS_PASSWORD")
	if !ok {
		output.Error(NoPasswordErr.Error())
		os.Exit(1)
	}
	return Credentials{username, password}, nil
}

func getURL() (string, error) {
	creds, err := getUserCreds()
	if err != nil {
		return "", err
	}
	var link = fmt.Sprintf("http://%s:%s@localhost:4200", creds.Username, creds.Password) // TODO: load link from file

	return link, nil
}

func getURLWithoutCreds() (string, error) {
	var link = fmt.Sprintf("http://localhost:4200") // TODO: load link from file
	return link, nil
}
