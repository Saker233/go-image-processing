package util

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)


func CreateHashed(pass []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword(pass, 10)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	return hashed, nil
}