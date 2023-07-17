package manager

import (
	"errors"

	"github.com/ChristianSassine/password-manager/pass-cli/internal/output"
)

var isQuiet = false

var (
	ServerContactErr   = errors.New("an error occured while contacting the server.")
	BadAuth            = errors.New("wrong password or username.")
	ServerErr          = errors.New("a server error occured.")
	ClientPasswordErr  = errors.New("bad request. The format of the request doesn't work with the server's API.")
	ClientClipboardErr = errors.New("unable to copy the password to the clipboard.")
)

func outputFormattedError(prefix string, err error) {
	output.Error("%v %v", prefix, err)
}
