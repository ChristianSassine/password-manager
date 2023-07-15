package manager

import (
	"errors"
	"net/http"
	"os"
	"regexp"
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
)

func CreateAccount(concatCreds string) {
	match, err := regexp.MatchString(credsPattern, concatCreds)
	if err != nil || !match {
		output.Error("%v %v", createAccountErrPrefix, ServerContactErr)
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
		output.Error("%v %v", createAccountErrPrefix, UserExistsErr)
		os.Exit(1)
	}

	output.Error("%v %v", createAccountErrPrefix, ServerContactErr)
	os.Exit(1)
}
