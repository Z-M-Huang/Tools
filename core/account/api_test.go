package account

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignupSuccess(t *testing.T) {
	request := &CreateAccountRequest{
		Email:           "test@example.com",
		Username:        "testUsername",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}

	status, response, tokenStr, expiresAt := SignUp(request)

	assert.Empty(t, response.Message)
	assert.NotEmpty(t, tokenStr)
	assert.Equal(t, http.StatusOK, status)
	assert.False(t, expiresAt.IsZero())
}

func TestSignupFail(t *testing.T) {
	requests := []*CreateAccountRequest{
		//Invalid Email
		{
			Email:           "adfadw",
			Username:        "testUsername",
			Password:        "abcdef123456",
			ConfirmPassword: "abcdef123456",
		},
		//Empty Email
		{
			Email:           "",
			Username:        "testUsername",
			Password:        "abcdef123456",
			ConfirmPassword: "abcdef123456",
		},
		//Password mismatch
		{
			Email:           "test@example.com",
			Username:        "testUsername",
			Password:        "abcdefabcdef",
			ConfirmPassword: "abcdef123456",
		},
		//Password length violation
		{
			Email:           "test@example.com",
			Username:        "testUsername",
			Password:        "abcdefaf",
			ConfirmPassword: "abcdefaf",
		},
	}

	for _, r := range requests {
		status, response, tokenStr, expiresAt := SignUp(r)
		assert.Empty(t, tokenStr)
		assert.NotEmpty(t, response.Message)
		assert.NotEqual(t, http.StatusOK, status)
		assert.NotEqual(t, http.StatusInternalServerError, status)
		assert.True(t, expiresAt.IsZero())
	}
}

func TestLoginSuccess(t *testing.T) {
	signupRequest := &CreateAccountRequest{
		Email:           "test@example.com",
		Username:        "testUsername",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	SignUp(signupRequest)
	loginRequest := &LoginRequest{
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	}

	status, response, tokenStr, expiresAt := login(loginRequest)

	assert.Equal(t, http.StatusOK, status)
	assert.Empty(t, response.Message)
	assert.NotEmpty(t, tokenStr)
	assert.False(t, expiresAt.IsZero())
	assert.True(t, expiresAt.After(time.Now()))
}

func TestLoginFail(t *testing.T) {
	signupRequest := &CreateAccountRequest{
		Email:           "test@example.com",
		Username:        "testUsername",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	SignUp(signupRequest)
	loginRequests := []*LoginRequest{
		//UserNotFound
		{
			Email:    "loginfail@test.com",
			Password: "abcdef123456",
		},
		//Password doesn't match
		{
			Email:    "test@example.com",
			Password: "abcdef1",
		},
	}
	for _, r := range loginRequests {
		status, response, tokenStr, expiresAt := login(r)
		assert.Empty(t, tokenStr)
		assert.NotEmpty(t, response.Message)
		assert.NotEqual(t, http.StatusOK, status)
		assert.NotEqual(t, http.StatusInternalServerError, status)
		assert.True(t, expiresAt.IsZero())
	}
}

func TestUpdatePasswordSuccess(t *testing.T) {
	signupRequest := &CreateAccountRequest{
		Email:           "updatePassword@example.com",
		Username:        "updatePassword",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	SignUp(signupRequest)
	request := &UpdatePasswordRequest{
		CurrentPassword: signupRequest.Password,
		Password:        "123456abcdef",
		ConfirmPassword: "123456abcdef",
	}

	status, response := updatePassword(request, signupRequest.Email)

	assert.Equal(t, http.StatusOK, status)
	assert.Empty(t, response.Message)
}

func TestUpdatePasswordFail(t *testing.T) {
	signupRequest := &CreateAccountRequest{
		Email:           "updatePasswordFail@example.com",
		Username:        "updatePasswordFail",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	SignUp(signupRequest)
	requests := []*UpdatePasswordRequest{
		//Request password doesn't match
		{
			CurrentPassword: signupRequest.Password,
			Password:        "abcdef123456",
			ConfirmPassword: "123456abcdef",
		},
		//Password minimum length
		{
			CurrentPassword: signupRequest.Password,
			Password:        "abcdef",
			ConfirmPassword: "abcdef",
		},
		//Password mismatch between database
		{
			CurrentPassword: "123456abcdef",
			Password:        "123456abcdef",
			ConfirmPassword: "123456abcdef",
		},
		//Password exact same
		{
			CurrentPassword: signupRequest.Password,
			Password:        signupRequest.Password,
			ConfirmPassword: signupRequest.Password,
		},
	}

	for _, r := range requests {
		status, response := updatePassword(r, signupRequest.Email)

		assert.NotEqual(t, http.StatusOK, status)
		assert.NotEqual(t, http.StatusInternalServerError, status)
		assert.NotEmpty(t, response.Message)
	}

	request := &UpdatePasswordRequest{
		CurrentPassword: signupRequest.Password,
		Password:        "123456abcdef",
		ConfirmPassword: "123456abcdef",
	}

	//User not found
	status, response := updatePassword(request, "notfound@test.com")
	assert.NotEqual(t, http.StatusOK, status)
	assert.NotEqual(t, http.StatusInternalServerError, status)
	assert.NotEmpty(t, response.Message)
}
