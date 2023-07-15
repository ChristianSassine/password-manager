package requests

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
)

func getURL() (string, error) {
	creds, err := getUserCreds()
	if err != nil {
		return "", err
	}
	url := mustGetURL()
	var link = fmt.Sprintf("http://%s:%s@%s", creds.Username, creds.Password, url)

	return link, nil
}

func getURLWithoutCreds() (string, error) {
	url := mustGetURL()
	var link = fmt.Sprintf("http://%s", url)
	return link, nil
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
