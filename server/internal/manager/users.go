package manager

import (
	"errors"

	"github.com/ChristianSassine/password-manager/server/internal/hashing"
	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	userAlreadyExistsErr = errors.New("user with the same username already exits.")
	userInvalidErr       = errors.New("Invalid credentials for user authentification.")
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

func validateUserCreds(username string, password string) error {
	filter := bson.D{{Key: usernameKey, Value: username}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}

	if !exists {
		return userInvalidErr
	}
	filter = bson.D{{Key: usernameKey, Value: username}}
	projection := bson.D{
		{Key: passwordKey, Value: 1},
	}
	res := mongodb.Get(filter, options.FindOne().SetProjection(projection))

	user := userData{}
	err = res.Decode(&user)
	if err != nil {
		return err
	}

	if user.Password != password {
		return userInvalidErr
	}

	return nil
}

func validateUserCreation(username string) error {
	filter := bson.D{{Key: usernameKey, Value: username}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}
	if exists {
		return userAlreadyExistsErr
	}
	return nil
}
