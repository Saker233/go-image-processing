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

func CompareHashed(hash, original []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(original))
	if err != nil {
		return err
	}
	return nil
}
