package usersvc

import (
	"net/mail"
	"slices"
	"strings"
	"unicode"

	"github.com/aedobrynin/soa-hw/core/internal/service"
)

func isASCIILetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

var (
	validPasswordSymbols = []rune{
		'_',
		'!',
		'@',
		'#',
		'$',
		'%',
		';',
		':',
		'^',
		'?',
		'*',
		'(',
		')',
		'-',
		'+',
		'=',
		'.',
		',',
	}
)

func isBadRuneForPassword(r rune) bool {
	return !(isASCIILetter(r) || unicode.IsDigit(r) || slices.Contains(validPasswordSymbols, r))
}

func validatePassword(password string) error {
	if len(password) < 10 || len(password) > 255 {
		return service.ErrPasswordValidation
	}
	if strings.ContainsFunc(password, isBadRuneForPassword) {
		return service.ErrPasswordValidation
	}
	return nil
}

func isBadRuneForLogin(r rune) bool {
	return !(isASCIILetter(r) || unicode.IsDigit(r) || r == '_')
}

func validateLogin(login string) error {
	if len(login) < 1 || len(login) > 25 {
		return service.ErrLoginValidation
	}
	if strings.ContainsFunc(login, isBadRuneForLogin) {
		return service.ErrLoginValidation
	}
	return nil
}

var isBadRuneForName = isBadRuneForLogin

func validateName(login string) error {
	if len(login) < 1 || len(login) > 25 {
		return service.ErrNameValidation
	}
	if strings.ContainsFunc(login, isBadRuneForName) {
		return service.ErrNameValidation
	}
	return nil
}

var isBadRuneForSurname = isBadRuneForLogin

func validateSurname(login string) error {
	if len(login) < 1 || len(login) > 25 {
		return service.ErrSurnameValidation
	}
	if strings.ContainsFunc(login, isBadRuneForSurname) {
		return service.ErrSurnameValidation
	}
	return nil
}

func validateEmail(email string) error {
	// TODO: length limits
	_, err := mail.ParseAddress(email)
	if err != nil {
		return service.ErrEmailValidation
	}
	return nil
}

func isBadRuneForPhone(r rune) bool {
	return !unicode.IsDigit(r)
}

func validatePhone(phone string) error {
	// TODO: support '+' symbol

	if len(phone) < 1 || len(phone) > 25 {
		return service.ErrPhoneValidation
	}

	if strings.ContainsFunc(phone, isBadRuneForPhone) {
		return service.ErrPhoneValidation
	}
	return nil
}
