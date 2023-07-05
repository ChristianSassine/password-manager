package manager

import (
	"errors"

	"github.com/ChristianSassine/password-manager/server/internal/hashing"
	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	usernameKey         = "username"
	passwordKey         = "password"
	managedPasswordsKey = "passwords"
)

type userData struct {
	Username  string            `bson:"username"`
	Password  string            `bson:"password"`
	Passwords map[string]string `bson:"passwords,omitempty"`
}

func CreateUser(username string, password string) error {
	err := validateUserCreation(username)
	if err != nil {
		return err
	}

	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = mongodb.Add(userData{Username: username, Password: hashedPassword})
	return err
}

func validateUserCreation(username string) error {
	filter := bson.D{{Key: usernameKey, Value: username}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Already Exists") // TODO: Create new errors
	}
	return nil
}
