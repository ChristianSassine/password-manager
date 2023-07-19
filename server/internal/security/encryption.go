package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"log"
	"os"
)

var NoKeyErr = errors.New("No secret key found for the encryption")

// Produces a 256 bit hashed key from the original key
func hashKey(key string) []byte {
	sha256Hash := sha256.Sum256([]byte(key))
	return []byte(sha256Hash[:])
}

// Encryption/Decryption using AES Block cipher and Galois/Counter Mode
func Encrypt(data []byte) ([]byte, error) {
	key, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Fatal(NoKeyErr)
	}
	hashedKey := hashKey(key)

	cipherBlock, err := aes.NewCipher(hashedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func Decrypt(cipherData []byte) ([]byte, error) {
	key, ok := os.LookupEnv("SECRET_KEY")
	if !ok {
		log.Fatal(NoKeyErr)
	}
	hashedKey := hashKey(key)

	cipherBlock, err := aes.NewCipher(hashedKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	nonce := cipherData[:gcm.NonceSize()]
	encryptedData := cipherData[gcm.NonceSize():]

	data, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
