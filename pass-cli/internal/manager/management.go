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
	PasswordExistsErr = errors.New("the key already has a password, please choose another key or remove the current one.")
	NoPasswordErr     = errors.New("the key doesn't have a password, please create the password before trying to access it.")
)

const (
	addErrPrefix    = "Adding password error:"
	getErrPrefix    = "Fetching password error:"
	renameErrPrefix = "Renaming password key error:"
	removeErrPrefix = "Renaming password key error:"
)

func AddPassword(opts requests.Parameters, returnPassword bool) {
	output.Print("Adding password...")

	response, err := requests.AddPassword(opts)
	if err != nil {
		outputFormattedError(addErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusCreated {
		output.Success("The password has been added!")
		GetPassword(opts.Key, returnPassword)
		return
	}

	if response.StatusCode == http.StatusForbidden {
		outputFormattedError(addErrPrefix, BadAuth)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusConflict {
		outputFormattedError(addErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusBadRequest {
		outputFormattedError(addErrPrefix, ServerErr)
		os.Exit(1)
	}

	outputFormattedError(addErrPrefix, ServerErr)
}

func GetPassword(key string, returnPassword bool) {
	output.Print("Fetching the password...")

	response, err := requests.GetPassword(key)
	if err != nil {
		outputFormattedError(getErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			outputFormattedError(getErrPrefix, err)
			os.Exit(1)
		}

		output.Success("The password has been fetched!")
		err = clipboard.WriteAll(string(body))
		if err != nil {
			outputFormattedError(getErrPrefix, ClientClipboardErr)
			os.Exit(1)
		}

		output.Print("It has been copied to the clipboard.")
		if returnPassword {
			output.AlwaysPrint(string(body))
		}
		return
	}

	if response.StatusCode == http.StatusForbidden {
		outputFormattedError(getErrPrefix, BadAuth)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusNotFound {
		outputFormattedError(getErrPrefix, NoPasswordErr)
		os.Exit(1)
	}
	outputFormattedError(getErrPrefix, ServerErr)
}

func RenameKey(key string, newKey string) {
	output.Print("Changing the password key...")

	response, err := requests.ChangePasswordKey(key, newKey)
	if err != nil {
		outputFormattedError(renameErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		output.Success("The key %s has been renamed to %s", key, newKey)
		return
	}

	if response.StatusCode == http.StatusForbidden {
		outputFormattedError(renameErrPrefix, BadAuth)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusConflict {
		outputFormattedError(renameErrPrefix, PasswordExistsErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusNotFound {
		outputFormattedError(renameErrPrefix, NoPasswordErr)
		os.Exit(1)
	}

	outputFormattedError(renameErrPrefix, ServerErr)
}

func RemovePassword(key string) {
	output.Print("Removing the password...")

	response, err := requests.RemovePassword(key)
	if err != nil {
		outputFormattedError(removeErrPrefix, ServerContactErr)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusOK {
		output.Success("The key %s has been removed.", key)
		return
	}

	if response.StatusCode == http.StatusForbidden {
		outputFormattedError(removeErrPrefix, BadAuth)
		os.Exit(1)
	}

	if response.StatusCode == http.StatusNotFound {
		outputFormattedError(removeErrPrefix, NoPasswordErr)
		os.Exit(1)
	}

	outputFormattedError(removeErrPrefix, ServerErr)
}
