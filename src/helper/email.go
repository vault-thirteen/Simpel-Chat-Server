package helper

import "net/mail"

func IsEmailAddressValid(email string) (isValid bool) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
