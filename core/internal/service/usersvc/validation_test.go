package usersvc

import (
	"strings"
	"testing"

	"github.com/aedobrynin/soa-hw/core/internal/service"
	"github.com/stretchr/testify/require"
)

func TestLoginValidation(t *testing.T) {
	testCases := []struct {
		Login   string
		IsValid bool
	}{
		{"", false},                      // Too short
		{strings.Repeat("a", 26), false}, // Too long
		{"русские буквы", false},         // Bad letters
		{"okokok@", false},               // Symbol
		{"validLogin", true},
		{"123213", true},
		{"val1dl0g1n", true},
		{"underscore_", true},
	}
	for _, tt := range testCases {
		t.Run(tt.Login, func(t *testing.T) {
			err := validateLogin(tt.Login)
			if tt.IsValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, service.ErrLoginValidation)
			}
		})
	}
}

func TestPasswordValidation(t *testing.T) {
	testCases := []struct {
		Password string
		IsValid  bool
	}{
		{strings.Repeat("a", 9), false},   // Too short
		{strings.Repeat("a", 256), false}, // Too long
		{"русские буквы", false},          // Bad letters
		{"_!@#$%;:^?*()-+=@.,", true},     // Symbols
		{"validPassword", true},
		{"1234567890", true},
		{"val1dpassw0rd", true},
		{"underscore_123", true},
	}
	for _, tt := range testCases {
		t.Run(tt.Password, func(t *testing.T) {
			err := validatePassword(tt.Password)
			if tt.IsValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, service.ErrPasswordValidation)
			}
		})
	}
}

func TestNameValidation(t *testing.T) {
	testCases := []struct {
		Name    string
		IsValid bool
	}{
		{"", false},                      // Too short
		{strings.Repeat("a", 26), false}, // Too long
		{"русские буквы", false},         // Bad letters
		{"_!@#$%;:^?*()-+=@.,", false},   // Symbols
		{"validName", true},
		{"1234567890", true},
		{"val1dname", true},
		{"underscore_123", true},
	}
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := validateName(tt.Name)
			if tt.IsValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, service.ErrNameValidation)
			}
		})
	}
}

func TestSurnameValidation(t *testing.T) {
	testCases := []struct {
		Surname string
		IsValid bool
	}{
		{"", false},                      // Too short
		{strings.Repeat("a", 26), false}, // Too long
		{"русские буквы", false},         // Bad letters
		{"_!@#$%;:^?*()-+=@.,", false},   // Symbols
		{"validName", true},
		{"1234567890", true},
		{"val1dname", true},
		{"underscore_123", true},
	}
	for _, tt := range testCases {
		t.Run(tt.Surname, func(t *testing.T) {
			err := validateSurname(tt.Surname)
			if tt.IsValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, service.ErrSurnameValidation)
			}
		})
	}
}

func TestEmailValidation(t *testing.T) {
	testCases := []struct {
		Email   string
		IsValid bool
	}{
		{"", false},
		{"русские буквы@mail.com", false}, // Bad letters
		{"_!@#$%;:^?*()-+=@.,", false},    // Symbols
		{"@", false},
		{"a@b", true}, // TODO: myb it's invalid?
		{"mail", false},
		{"mail@", false},
		{"@mail.com", false},
		{"mail@mail.com", true},
		{"my_email@my_email.bz", true},
	}
	for _, tt := range testCases {
		t.Run(tt.Email, func(t *testing.T) {
			err := validateEmail(tt.Email)
			if tt.IsValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, service.ErrEmailValidation)
			}
		})
	}
}

func TestPhoneValidation(t *testing.T) {
	testCases := []struct {
		Phone   string
		IsValid bool
	}{
		{"", false},                      // Too short
		{strings.Repeat("1", 26), false}, // Too long
		{"letters123", false},            // Letters
		{"_12312312", false},             // Symbols
		{"123213123", true},
		{"1", true},
		{"555", true},
	}
	for _, tt := range testCases {
		t.Run(tt.Phone, func(t *testing.T) {
			err := validatePhone(tt.Phone)
			if tt.IsValid {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, service.ErrPhoneValidation)
			}
		})
	}
}
