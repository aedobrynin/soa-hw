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
	// TODO
}

func TestSurnameValidation(t *testing.T) {
	// TODO
}

func TestEmailValidation(t *testing.T) {
	// TODO
}

func TestPhoneValidation(t *testing.T) {
	// TODO
}
