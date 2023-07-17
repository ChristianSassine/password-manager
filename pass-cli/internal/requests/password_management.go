package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
)

const (
	UserUsernameEnv = "PASS_USERNAME"
	UserPasswordEnv = "PASS_PASSWORD"
)

const (
	passwordPath = "password"
	passKeyQuery = "key"
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

type Parameters struct {
	Key          string `json:"key"`
	Length       int    `json:"length"`
	LowerLetters bool   `json:"lowerLetters"`
	UpperLetters bool   `json:"upperLetters"`
	Digits       bool   `json:"digits"`
	Symbols      bool   `json:"symbols"`
}

var client = &http.Client{}

func GetPassword(key string) (*http.Response, error) {
	url, err := getURL(passwordPath, query{passKeyQuery, key})
	if err != nil {
		return nil, err
	}
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AddPassword(opts Parameters) (*http.Response, error) {
	url, err := getURL(passwordPath)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	response, err := client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func ChangePasswordKey(oldKey string, newKey string) (*http.Response, error) {
	url, err := getURL(passwordPath)
	if err != nil {
		return nil, err
	}

	change := &keysChange{OldKey: oldKey, NewKey: newKey}
	data, err := json.Marshal(change)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
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
	url, err := getURL(passwordPath, query{passKeyQuery, key})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, url, &bytes.Buffer{})
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
	username, ok := os.LookupEnv(UserUsernameEnv)
	if !ok {
		output.Error(NoUsernameErr.Error())
		os.Exit(1)
	}

	password, ok := os.LookupEnv(UserPasswordEnv)
	if !ok {
		output.Error(NoPasswordErr.Error())
		os.Exit(1)
	}
	return Credentials{username, password}, nil
}
