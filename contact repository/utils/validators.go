package utils

import (
	"contactapp/models/user"
	"errors"
	"net/mail"
)

func UserValidator(user user.User) error {
	if len(user.FullName) < 5 {
		return errors.New("fullname cannot be less than 5")
	}

	if len(user.Password) < 7 {
		return errors.New("password cannot be less than 7")
	}

	if !isEmailValid(user.Email) {
		return errors.New("Invalid Email")
	}

	return nil
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}