package password

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func RandomPassword() string {
	var digitBytes string = "0123456789"
	
	password := make([]byte, 4)
	for i := range password {
		password[i] = digitBytes[rand.Intn(len(digitBytes))]
	}

	return string(password)
}