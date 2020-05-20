package account

import (
	"net/http"
	"testing"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func init() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	data.Config = &data.Configuration{
		DatabaseConfig: &data.DatabaseConfiguration{
			ConnectionString: "./test.db",
			Driver:           "sqlite3",
		},
		RedisConfig: &data.RedisConfiguration{
			Addr: mr.Addr(),
		},
		GoogleOauthConfig: &data.GoogleOauthConfiguration{
			ClientID:     "testClientID",
			ClientSecret: "testClientSecret",
		},
		JwtKey:          []byte("CBYtDWTfRU5Pv7yULj46vm8ueZG7hbnq"),
		Host:            "localhost",
		ResourceVersion: "1",
		IsDebug:         true,
		HTTPS:           false,
		EnableSitemap:   true,
	}
}

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	data.Config = &data.Configuration{
		DatabaseConfig: &data.DatabaseConfiguration{
			ConnectionString: "./test.db",
			Driver:           "sqlite3",
		},
		RedisConfig: &data.RedisConfiguration{
			Addr: mr.Addr(),
		},
		GoogleOauthConfig: &data.GoogleOauthConfiguration{
			ClientID:     "testClientID",
			ClientSecret: "testClientSecret",
		},
		JwtKey:          []byte("CBYtDWTfRU5Pv7yULj46vm8ueZG7hbnq"),
		Host:            "localhost",
		ResourceVersion: "1",
		IsDebug:         true,
		HTTPS:           false,
		EnableSitemap:   true,
	}
	m.Run()
}

func TestSignupSuccess(t *testing.T) {
	request := &CreateAccountRequest{
		Email:           "test@example.com",
		Username:        "testUsername",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}

	status, response, tokenStr, expiresAt := signUp(request)

	assert.Equal(t, "", response.ErrorMessage, "Error Message is not empty")
	assert.NotEqual(t, "", tokenStr, "Token is empty")
	assert.Equal(t, http.StatusOK, status, "Status is not 200")
	assert.False(t, expiresAt.IsZero(), "ExpiresAt is zero time")
}

func TestSignupFail(t *testing.T) {
	t.Error("?")
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
		status, response, tokenStr, expiresAt := signUp(r)
		t.Log(response.ErrorMessage)
		assert.Equal(t, "", tokenStr, "Token is not empty in signup fail")
		assert.Equal(t, "", response.ErrorMessage, "ErrorMessage in signup is empty on fail")
		assert.NotEqual(t, http.StatusOK, status, "Signup status is 200 on fail")
		assert.True(t, expiresAt.IsZero(), "Signup expiresAt is not zero time on fail")
	}
}

func TestT(t *testing.T) {
	assert.True(t, false, "aha")
}
