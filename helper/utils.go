package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// ComparePassword is used to compare the password.
func ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
