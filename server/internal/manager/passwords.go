package manager

import (
	"errors"

	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
	"github.com/ChristianSassine/password-manager/server/internal/security"
	"github.com/ChristianSassine/password-manager/server/pkg/generator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	NoPasswordErr       = errors.New("Invalid password key. The password doesn't exist")
	PasswordConflictErr = errors.New("Invalid password key. The password already exists")
)

type PasswordOptions struct {
	Key string `json:"key"`
	generator.Options
}

func UserAddPassword(username string, opts PasswordOptions) error {
	var filter = bson.D{{Key: usernameKey, Value: username}, {Key: managedPasswordsKey + "." + opts.Key, Value: bson.D{{Key: "$exists", Value: true}}}}
	exists, err := mongodb.Exist(filter)
	if err != nil {
		return err
	}
	if exists {
		return PasswordConflictErr
	}

	return addPassword(username, opts.Key, opts.Options)
}

func UserGetPassword(username string, key string) (string, error) {
	return getPassword(username, key)
}

func UserRemovePassword(username string, key string) error {
	return removePassword(username, key)
}

func UserRenamePassword(username string, oldKey string, newKey string) error {
	return renamePasswordKey(username, oldKey, newKey)
}

func addPassword(username string, key string, opts generator.Options) error {
	var filter = bson.D{{Key: usernameKey, Value: username}}

	newPassword, err := generator.Generate(opts)
	if err != nil {
		return err
	}
	encryptedPass, err := security.Encrypt([]byte(newPassword))
	if err != nil {
		return err
	}

	var update = bson.D{{Key: "$set", Value: bson.D{{Key: managedPasswordsKey + "." + key, Value: encryptedPass}}}}
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
	decryptedPass, err := security.Decrypt(data.Passwords[key])
	if err != nil {
		return "", err
	}

	return string(decryptedPass), nil
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

	filter = bson.D{{Key: usernameKey, Value: username}, {Key: managedPasswordsKey + "." + newKey, Value: bson.D{{Key: "$exists", Value: true}}}}
	exists, err = mongodb.Exist(filter)
	if err != nil {
		return err
	}
	if exists {
		return PasswordConflictErr
	}

	filter = bson.D{{Key: usernameKey, Value: username}}
	var update = bson.D{{Key: "$rename", Value: bson.D{{Key: managedPasswordsKey + "." + oldKey, Value: managedPasswordsKey + "." + newKey}}}}
	_, err = mongodb.Update(filter, update)
	if err != nil {
		return err
	}

	return nil
}
