package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)

}

