package manager

import (
	"errors"

	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
	"github.com/ChristianSassine/password-manager/server/pkg/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	NoPasswordErr       = errors.New("Invalid password key. The password doesn't exist")
	PasswordConflictErr = errors.New("Invalid password key. The password already exists")
)

// TODO: Migrate to using mongoDB ids instead of usernames

func UserAddPassword(username string, userPassword string, key string) error {
	var filter = bson.D{{Key: usernameKey, Value: username}, {Key: managedPasswordsKey + "." + key, Value: bson.D{{Key: "$exists", Value: true}}}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}
	if exists {
		return PasswordConflictErr
	}

	if err := validateUserCreds(username, userPassword); err != nil {
		return err
	}

	return addPassword(username, key)
}

func UserGetPassword(username string, userPassword string, key string) (string, error) {
	if err := validateUserCreds(username, userPassword); err != nil {
		return "", err
	}

	return getPassword(username, key)
}

func UserRemovePassword(username string, userPassword string, key string) error {
	if err := validateUserCreds(username, userPassword); err != nil {
		return err
	}

	return removePassword(username, key)
}

func UserRenamePassword(username string, userPassword string, oldKey string, newKey string) error {
	if err := validateUserCreds(username, userPassword); err != nil {
		return err
	}

	return renamePasswordKey(username, oldKey, newKey)
}

// TODO: Add password encryption
func addPassword(username string, key string) error {
	var filter = bson.D{{Key: usernameKey, Value: username}}

	newPassword, err := password.Generate(password.Options{Length: 32, LowerLetters: true, UpperLetters: true, Digits: true, Symbols: true}) // TODO: Make password configurable
	if err != nil {
		return err
	}
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: managedPasswordsKey + "." + key, Value: newPassword}}}}
	_, err = mongodb.Update(filter, update)
	if err != nil {
		return err
	}
	return nil
}

func getPassword(username string, key string) (string, error) {
	var filter = bson.D{{Key: usernameKey, Value: username}, {Key: managedPasswordsKey + "." + key, Value: bson.D{{Key: "$exists", Value: true}}}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", NoPasswordErr
	}

	opts := options.FindOne().SetProjection(bson.D{{Key: managedPasswordsKey + "." + key, Value: 1}})
	filter = bson.D{{Key: usernameKey, Value: username}, {Key: managedPasswordsKey + "." + key, Value: bson.D{{Key: "$exists", Value: true}}}}
	res := mongodb.Get(filter, opts)

	var data UserData
	if err := res.Decode(&data); err != nil {
		return "", err
	}
	return data.Passwords[key], nil
}

func removePassword(username string, key string) error {
	var filter = bson.D{{Key: usernameKey, Value: username}}
	var update = bson.D{{Key: "$unset", Value: bson.D{{Key: managedPasswordsKey + "." + key}}}}
	_, err := mongodb.Update(filter, update)
	if err != nil {
		return err
	}
	return nil
}

func renamePasswordKey(username string, oldKey string, newKey string) error {
	var filter = bson.D{{Key: usernameKey, Value: username}, {Key: managedPasswordsKey + "." + oldKey, Value: bson.D{{Key: "$exists", Value: true}}}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}
	if !exists {
		return NoPasswordErr
	}
	filter = bson.D{{Key: usernameKey, Value: username}}
	var update = bson.D{{Key: "$rename", Value: bson.D{{Key: managedPasswordsKey + "." + oldKey, Value: managedPasswordsKey + "." + newKey}}}}
	_, err = mongodb.Update(filter, update)
	if err != nil {
		return err
	}
	return nil
}
