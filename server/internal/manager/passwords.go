package manager

import (
	"github.com/ChristianSassine/password-manager/server/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func addPassword(username string, key string, password string) error {
	var filter = bson.D{{Key: usernameKey, Value: username}}
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: managedPasswordsKey + "." + key, Value: password}}}}
	_, err := mongodb.Update(filter, update)
	if err != nil {
		return err
	}
	return nil
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

func renameKey(username string, oldKey string, newKey string) error {
	var filter = bson.D{{Key: usernameKey, Value: username}}
	var update = bson.D{{Key: "$rename", Value: bson.D{{Key: managedPasswordsKey + "." + oldKey, Value: managedPasswordsKey + "." + newKey}}}}
	_, err := mongodb.Update(filter, update)
	if err != nil {
		return err
	}
	return nil
}
