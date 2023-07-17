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

func getURL(queries ...query) (string, error) {
	creds, err := getUserCreds()
	if err != nil {
		return "", err
	}

	urlPath := mustGetURL()
	var link = fmt.Sprintf("http://%s:%s@%s", creds.Username, creds.Password, urlPath)
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	for _, q := range queries {
		u.Query().Set(q.key, q.value)
	}

	return u.String(), nil
}

func getURLWithoutCreds() (string, error) {
	urlPath := mustGetURL()
	var link = fmt.Sprintf("http://%s", urlPath)
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
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

	output.Success(path)
	return string(b)
}
