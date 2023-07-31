package manager

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/requests"
)

const credsPattern = "^[^:\n]+:[^:\n]+$"

var (
	CredsPatternErr = errors.New("invalid username and password. Please follow the pattern <username>:<password>")
	UserExistsErr   = errors.New("user already exists. Please choose another username")
)

const (
	createAccountErrPrefix = "Creating account error:"
	setURLErrPrefix        = "Updating URL error :"
	gettingURLErrPrefix    = "Getting URL error :"
)

func CreateAccount(concatCreds string) {
	match, err := regexp.MatchString(credsPattern, concatCreds)
	if err != nil || !match {
		outputFormattedError(createAccountErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	splitCreds := strings.Split(concatCreds, ":")
	creds := requests.Credentials{Username: splitCreds[0], Password: splitCreds[1]}

	resp, err := requests.CreateAccountRequest(creds)
	if resp.StatusCode == http.StatusCreated {
		output.Success("The account has been created! You can now change the environment variables.")
		return
	}

	if resp.StatusCode == http.StatusConflict {
		outputFormattedError(createAccountErrPrefix, UserExistsErr)
		os.Exit(1)
	}

	outputFormattedError(createAccountErrPrefix, ServerContactErr)
	os.Exit(1)
}

func SetURL(rawURL string) {
	OS := runtime.GOOS
	configPath, err := os.UserConfigDir()
	if err != nil {
		outputFormattedError(setURLErrPrefix, err)
		os.Exit(1)
	}

	var path string
	if OS == "windows" {
		configPath += "\\pass-cli"
		path = fmt.Sprintf("%s\\server.txt", configPath)
	} else {
		configPath += "/pass-cli"
		path = fmt.Sprintf("%s/server.txt", configPath)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		outputFormattedError(setURLErrPrefix, err)
		os.Exit(1)
	}

	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		outputFormattedError(setURLErrPrefix, err)
		os.Exit(1)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		outputFormattedError(setURLErrPrefix, err)
		os.Exit(1)
	}

	_, err = f.WriteString(u.String())
	if err != nil {
		outputFormattedError(setURLErrPrefix, err)
		os.Exit(1)
	}

	output.Success("Success in setting the new URL!")
}
