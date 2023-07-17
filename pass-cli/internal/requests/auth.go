package requests

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

const (
	credsPattern = "^[^:\n]+:[^:\n]+$"
	authPath     = "/auth"
)

func CreateAccountRequest(creds Credentials) (*http.Response, error) {
	url, err := getURLWithoutCreds(authPath)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	return client.Post(url, "application/json", bytes.NewReader(b))
}
