package manager

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
	"github.com/ChristianSassine/password-manager/pass-cli/internal/requests"
	"github.com/atotto/clipboard"
)

var (
	ServerContactErr   = errors.New("an error occured while contacting the server.")
	ServerErr          = errors.New("a server error occured.")
	PasswordExistsErr  = errors.New("the key already has a password, please choose another key or remove the current one.")
	NoPasswordErr      = errors.New("the key doesn't have a password, please create the password before trying to access it.")
	ClientPasswordErr  = errors.New("bad request. The format of the request doesn't work with the server's API.")
	ClientClipboardErr = errors.New("unable to copy the password to the clipboard.")
)

const (
	AddErrPrefix = "Adding password error:"
	GetErrPrefix = "Fetching password error:"
)

func AddPassword(key string) {
	response, err := requests.AddPassword(key)
	output.NormalLn("Adding password...")

	if err != nil {
		output.Error("%v %v", AddErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusCreated {
		output.Success("The password has been added!")
		// TODO: Fetch password
		return
	}

	if response.StatusCode == http.StatusConflict {
		output.Error("%v %v", AddErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusBadRequest {
		output.Error("%v %v", AddErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}
	output.Error("%v %v", AddErrPrefix, ServerErr)
}

func GetPassword(key string) {
	response, err := requests.GetPassword(key)
	output.NormalLn("Fetching the password...")

	if err != nil {
		output.Error("%v %v", GetErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			output.Error("%v %v", GetErrPrefix, err)
			os.Exit(1)
		}

		output.Success("The password has been fetched!")
		err = clipboard.WriteAll(string(body))
		if err != nil {
			output.Error("%v %v", GetErrPrefix, ClientClipboardErr)
			os.Exit(1)
		}

		output.NormalLn("It has been copied to the clipboard.")
		return
	}

	if response.StatusCode == http.StatusNotFound {
		output.Error("%v %v", GetErrPrefix, NoPasswordErr)
		os.Exit(1)
	}
	output.Error("%v %v", GetErrPrefix, ServerErr)
}

func ChangePasswordKey(key string) {
	response, err := requests.GetPassword(key)
	output.NormalLn("Fetching the password...")

	if err != nil {
		output.Error("%v %v", GetErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			output.Error("%v %v", GetErrPrefix, err)
			os.Exit(1)
		}

		output.Success("The password has been fetched!")
		err = clipboard.WriteAll(string(body))
		if err != nil {
			output.Error("%v %v", GetErrPrefix, ClientClipboardErr)
			os.Exit(1)
		}

		output.NormalLn("It has been copied to the clipboard.")
		return
	}

	if response.StatusCode == http.StatusNotFound {
		output.Error("%v %v", GetErrPrefix, NoPasswordErr)
		os.Exit(1)
	}
	output.Error("%v %v", GetErrPrefix, ServerErr)
}
