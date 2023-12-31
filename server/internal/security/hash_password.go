package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bcryptCost := 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
