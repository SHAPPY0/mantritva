package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSaltPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println("Error while generating password hash", err.Error())
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashPassword string, password string) bool {
	byteHash := []byte(hashPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(password))
	if err != nil {
		fmt.Println("Exception while comparing password", err.Error())
		return false
	}
	return true
}
