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
	addErrPrefix    = "Adding password error:"
	getErrPrefix    = "Fetching password error:"
	renameErrPrefix = "Renaming password key error:"
	removeErrPrefix = "Renaming password key error:"
)

func AddPassword(key string) {
	output.NormalLn("Adding password...")

	response, err := requests.AddPassword(key)
	if err != nil {
		output.Error("%v %v", addErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusCreated {
		output.Success("The password has been added!")
		GetPassword(key)
		return
	}

	if response.StatusCode == http.StatusConflict {
		output.Error("%v %v", addErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusBadRequest {
		output.Error("%v %v", addErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}
	output.Error("%v %v", addErrPrefix, ServerErr)
}

func GetPassword(key string) {
	output.NormalLn("Fetching the password...")

	response, err := requests.GetPassword(key)
	if err != nil {
		output.Error("%v %v", getErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			output.Error("%v %v", getErrPrefix, err)
			os.Exit(1)
		}

		output.Success("The password has been fetched!")
		err = clipboard.WriteAll(string(body))
		if err != nil {
			output.Error("%v %v", getErrPrefix, ClientClipboardErr)
			os.Exit(1)
		}

		output.NormalLn("It has been copied to the clipboard.")
		return
	}

	if response.StatusCode == http.StatusNotFound {
		output.Error("%v %v", getErrPrefix, NoPasswordErr)
		os.Exit(1)
	}
	output.Error("%v %v", getErrPrefix, ServerErr)
}

func RenameKey(key string, newKey string) {
	output.NormalLn("Changing the password key...")

	response, err := requests.ChangePasswordKey(key, newKey)
	if err != nil {
		output.Error("%v %v", renameErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		output.Success("The key %s has been renamed to %s", key, newKey)
		return
	}

	if response.StatusCode == http.StatusConflict {
		output.Error("%v %v", renameErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusNotFound {
		output.Error("%v %v", renameErrPrefix, NoPasswordErr)
		os.Exit(1)
	}

	output.Error("%v %v", renameErrPrefix, ServerErr)
}

func RemovePassword(key string) {
	output.NormalLn("Removing the password...")

	response, err := requests.RemovePassword(key)
	if err != nil {
		output.Error("%v %v", removeErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		output.Success("The key %s has been removed.", key)
		return
	}

	if response.StatusCode == http.StatusNotFound {
		output.Error("%v %v", removeErrPrefix, NoPasswordErr)
		os.Exit(1)
	}

	output.Error("%v %v", removeErrPrefix, ServerErr)
}
