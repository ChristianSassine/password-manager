package requests

import (
	"fmt"
	"net/url"
	"os"
	"runtime"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
)

type query struct {
	key   string
	value string
}

func getURL(path string, queries ...query) (string, error) {
	creds, err := getUserCreds()
	if err != nil {
		return "", err
	}

	baseURL := mustGetURL()
	var link = fmt.Sprintf("http://%s:%s@%s", creds.Username, creds.Password, baseURL)
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	u = u.JoinPath(path)
	urlQuery := u.Query()
	for _, q := range queries {
		urlQuery.Set(q.key, q.value)
	}
	u.RawQuery = urlQuery.Encode()

	return u.String(), nil
}

func getURLWithoutCreds(path string) (string, error) {
	baseURL := mustGetURL()
	var link = fmt.Sprintf("http://%s", baseURL)
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	u = u.JoinPath(path)
	return u.String(), nil
}

func mustGetURL() string {
	OS := runtime.GOOS
	configPath, err := os.UserConfigDir()
	if err != nil {
		output.Error("Error while getting the URL: %s", err)
		os.Exit(1)
	}

	var path string
	if OS == "windows" {
		path = fmt.Sprintf("%s\\pass-cli\\server.txt", configPath)
	}

	if OS == "linux" {
		path = fmt.Sprintf("%s/pass-cli/server.txt", configPath)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		output.Error("Error while getting the URL: %s", err)
		os.Exit(1)
	}

	return string(b)
}
