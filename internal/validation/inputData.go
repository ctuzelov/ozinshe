package validation

import (
	"errors"
	"net/mail"
	"ozinshe/internal/models"
	"regexp"
)

const (
	MsgInvalidEmail = "write correct email"
	MsgInvalidPass  = "password must contain letters, numbers and must be at least 6 characters"
)

func GetErrMsg(m models.User) error {
	if !IsValidEmail(m.Email) {
		return errors.New(MsgInvalidEmail)
	}
	if !IsValidPassword(m.Password) {
		return errors.New(MsgInvalidPass)
	}
	return nil
}

func IsValidEmail(email string) bool {
	rxEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(pass string) bool {
	tests := []string{".{6,}", "[a-z]", "[0-9]"}
	for _, test := range tests {
		valid, _ := regexp.MatchString(test, pass)
		if !valid {
			return false
		}
	}
	return true
}
