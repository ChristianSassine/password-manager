package manager

import (
	"errors"

	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
	"github.com/ChristianSassine/password-manager/server/internal/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UserAlreadyExistsErr = errors.New("user with the same username already exits.")
	UserInvalidErr       = errors.New("Invalid credentials for user authentification.")
)

const (
	usernameKey         = "username"
	passwordKey         = "password"
	managedPasswordsKey = "passwords"
)

type UserData struct {
	Username  string            `bson:"username"`
	Password  string            `bson:"password"`
	Passwords map[string][]byte `bson:"passwords,omitempty"`
}

func CreateUser(username string, password string) error {
	err := validateUserCreation(username)
	if err != nil {
		return err
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = mongodb.Add(UserData{Username: username, Password: hashedPassword})
	return err
}

func ValidateUserCreds(username string, password string) error {
	filter := bson.D{{Key: usernameKey, Value: username}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}

	if !exists {
		return UserInvalidErr
	}
	filter = bson.D{{Key: usernameKey, Value: username}}
	projection := bson.D{
		{Key: passwordKey, Value: 1},
	}
	res := mongodb.Get(filter, options.FindOne().SetProjection(projection))

	user := UserData{}
	err = res.Decode(&user)
	if err != nil {
		return err
	}

	if !security.CheckPassword(password, user.Password) {
		return UserInvalidErr
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
		return UserAlreadyExistsErr
	}
	return nil
}
