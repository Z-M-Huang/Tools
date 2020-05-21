package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
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

	db.InitDB()
	db.InitRedis()
	account.InitGoogleOauth()
}

func teardown() {
	err := db.Disconnect()
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		err = os.Remove(data.Config.DatabaseConfig.ConnectionString)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPageStyle(t *testing.T) {
	styles := []*data.PageStyleData{
		{
			Name:      "Cerulean",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cerulean/bootstrap.min.css",
			Integrity: "sha384-LV/SIoc08vbV9CCeAwiz7RJZMI5YntsH8rGov0Y2nysmepqMWVvJqds6y0RaxIXT",
		},
		{
			Name:      "Cosmo",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cosmo/bootstrap.min.css",
			Integrity: "sha384-qdQEsAI45WFCO5QwXBelBe1rR9Nwiss4rGEqiszC+9olH1ScrLrMQr1KmDR964uZ",
		},
		{
			Name:      "Cyborg",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cyborg/bootstrap.min.css",
			Integrity: "sha384-l7xaoY0cJM4h9xh1RfazbgJVUZvdtyLWPueWNtLAphf/UbBgOVzqbOTogxPwYLHM",
		},
		{
			Name:      "Darkly",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/darkly/bootstrap.min.css",
			Integrity: "sha384-rCA2D+D9QXuP2TomtQwd+uP50EHjpafN+wruul0sXZzX/Da7Txn4tB9aLMZV4DZm",
		},
		{
			Name:      "Flatly",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/flatly/bootstrap.min.css",
			Integrity: "sha384-yrfSO0DBjS56u5M+SjWTyAHujrkiYVtRYh2dtB3yLQtUz3bodOeialO59u5lUCFF",
		},
		{
			Name:      "Journal",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/journal/bootstrap.min.css",
			Integrity: "sha384-0d2eTc91EqtDkt3Y50+J9rW3tCXJWw6/lDgf1QNHLlV0fadTyvcA120WPsSOlwga",
		},
		{
			Name:      "Litera",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/litera/bootstrap.min.css",
			Integrity: "sha384-pLgJ8jZ4aoPja/9zBSujjzs7QbkTKvKw1+zfKuumQF9U+TH3xv09UUsRI52fS+A6",
		},
		{
			Name:      "Lumen",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/lumen/bootstrap.min.css",
			Integrity: "sha384-715VMUUaOfZR3/rWXZJ9agJ/udwSXGEigabzUbJm2NR4/v5wpDy8c14yedZN6NL7",
		},
		{
			Name:      "Lux",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/lux/bootstrap.min.css",
			Integrity: "sha384-oOs/gFavzADqv3i5nCM+9CzXe3e5vXLXZ5LZ7PplpsWpTCufB7kqkTlC9FtZ5nJo",
		},
		{
			Name:      "Materia",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/materia/bootstrap.min.css",
			Integrity: "sha384-1tymk6x9Y5K+OF0tlmG2fDRcn67QGzBkiM3IgtJ3VrtGrIi5ryhHjKjeeS60f1FA",
		},
		{
			Name:      "Minty",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/minty/bootstrap.min.css",
			Integrity: "sha384-4HfFay3AYJnEmbgRzxYWJk/Ka5jIimhB/Fssk7NGT9Tj3rkEChpSxLK0btOGzf2I",
		},
		{
			Name:      "Pulse",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/pulse/bootstrap.min.css",
			Integrity: "sha384-FnujoHKLiA0lyWE/5kNhcd8lfMILbUAZFAT89u11OhZI7Gt135tk3bGYVBC2xmJ5",
		},
		{
			Name:      "Sandstone",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/sandstone/bootstrap.min.css",
			Integrity: "sha384-ABdnjefqVzESm+f9z9hcqx2cvwvDNjfrwfW5Le9138qHCMGlNmWawyn/tt4jR4ba",
		},
		{
			Name:      "Simplex",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/simplex/bootstrap.min.css",
			Integrity: "sha384-cRAmF0wErT4D9dEBc36qB6pVu+KmLh516IoGWD/Gfm6FicBbyDuHgS4jmkQB8u1a",
		},
		{
			Name:      "Sketchy",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/sketchy/bootstrap.min.css",
			Integrity: "sha384-2kOE+STGAkgemIkUbGtoZ8dJLqfvJ/xzRnimSkQN7viOfwRvWseF7lqcuNXmjwrL",
		},
		{
			Name:      "Salte",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/slate/bootstrap.min.css",
			Integrity: "sha384-G9YbB4o4U6WS4wCthMOpAeweY4gQJyyx0P3nZbEBHyz+AtNoeasfRChmek1C2iqV",
		},
		{
			Name:      "Solar",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/solar/bootstrap.min.css",
			Integrity: "sha384-Ya0fS7U2c07GI3XufLEwSQlqhpN0ri7w/ujyveyqoxTJ2rFHJcT6SUhwhL7GM4h9",
		},
		{
			Name:      "Spacelab",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/spacelab/bootstrap.min.css",
			Integrity: "sha384-nl8CRcNYOGaXP68QRJeVTXCWAhwqBhRha0QbuubX1qDQpGT3GgclpvyzkR2JzyfZ",
		},
		{
			Name:      "Superhero",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/superhero/bootstrap.min.css",
			Integrity: "sha384-R/oa7KS0iDoHwdh4Gyl3/fU7pgvSCt7oyuQ79pkw+e+bMWD9dzJJa+Zqd+XJS0AD",
		},
		{
			Name:      "United",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/united/bootstrap.min.css",
			Integrity: "sha384-bzjLLgZOhgXbSvSc5A9LWWo/mSIYf7U7nFbmYIB2Lgmuiw3vKGJuu+abKoaTx4W6",
		},
		{
			Name:      "Yeti",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/yeti/bootstrap.min.css",
			Integrity: "sha384-bWm7zrSUE5E+21rA9qdH5frkMpXvqjQm/WJw4L5PYQLNHrywI5zs5saEjIcCdGu1",
		},
		{
			Name:      "Default",
			Link:      "https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css",
			Integrity: "sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk",
		},
	}

	for _, s := range styles {
		ret := getPageStyle(strings.ToLower(s.Name))

		assert.Equal(t, s.Name, ret.Name)
		assert.Equal(t, s.Link, ret.Link)
		assert.Equal(t, s.Integrity, ret.Integrity)
	}
}

func TestApiAuthHandler(t *testing.T) {
	router := SetupRouter()
	signupRequest := &account.CreateAccountRequest{
		Email:           "testHomeAPIAuthHandler@example.com",
		Username:        "testHomeAPIAuthHandler",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	requestStr, err := json.Marshal(signupRequest)
	assert.Empty(t, err)
	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(requestStr))
	req.Header.Add("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authCookie *http.Cookie

	for _, c := range w.Result().Cookies() {
		if c.Name == utils.SessionCookieKey {
			authCookie = c
			break
		}
	}

	loginRequest := &account.LoginRequest{
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	}
	requestStr, err = json.Marshal(loginRequest)
	assert.Empty(t, err)
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(requestStr))
	req.Header.Add("content-type", "application/json")
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApiAuthHandlerFailed(t *testing.T) {
	router := SetupRouter()
	signupRequest := &account.CreateAccountRequest{
		Email:           "TestApiAuthHandlerFailed@example.com",
		Username:        "TestApiAuthHandlerFailed",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	requestStr, err := json.Marshal(signupRequest)
	assert.Empty(t, err)
	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(requestStr))
	req.Header.Add("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authCookie *http.Cookie

	for _, c := range w.Result().Cookies() {
		if c.Name == utils.SessionCookieKey {
			authCookie = c
			break
		}
	}

	assert.NotEmpty(t, authCookie)
	authCookie.Value = "asdfasdfasd12312"
	loginRequest := &account.UpdatePasswordRequest{
		CurrentPassword: signupRequest.Password,
		Password:        "123456abcdef",
		ConfirmPassword: "123456abcdef",
	}
	requestStr, err = json.Marshal(loginRequest)
	assert.Empty(t, err)
	req, _ = http.NewRequest("POST", "/api/account/update/password", bytes.NewBuffer(requestStr))
	req.Header.Add("content-type", "application/json")
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPageAuthHandler(t *testing.T) {
	router := SetupRouter()
	signupRequest := &account.CreateAccountRequest{
		Email:           "TestPageAuthHandler@example.com",
		Username:        "TestPageAuthHandler",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	requestStr, err := json.Marshal(signupRequest)
	assert.Empty(t, err)
	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(requestStr))
	req.Header.Add("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authCookie *http.Cookie

	for _, c := range w.Result().Cookies() {
		if c.Name == utils.SessionCookieKey {
			authCookie = c
			break
		}
	}

	req, _ = http.NewRequest("GET", "/", nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPageAuthHandlerFailed(t *testing.T) {
	router := SetupRouter()
	signupRequest := &account.CreateAccountRequest{
		Email:           "TestPageAuthHandlerFailed@example.com",
		Username:        "TestPageAuthHandlerFailed",
		Password:        "abcdef123456",
		ConfirmPassword: "abcdef123456",
	}
	requestStr, err := json.Marshal(signupRequest)
	assert.Empty(t, err)
	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(requestStr))
	req.Header.Add("content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authCookie *http.Cookie

	for _, c := range w.Result().Cookies() {
		if c.Name == utils.SessionCookieKey {
			authCookie = c
			break
		}
	}

	assert.NotEmpty(t, authCookie)
	authCookie.Value = "asdfasdfasd12312"
	req, _ = http.NewRequest("GET", "/account", nil)
	req.AddCookie(authCookie)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
}

func TestPageStyleHandler(t *testing.T) {
	router := SetupRouter()

	c := &http.Cookie{
		Name:   utils.PageStyleCookieKey,
		Value:  "default",
		Domain: data.Config.Host,
	}
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(c)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("content-type"), "text/html")
}

func TestHomePage(t *testing.T) {
	router := SetupRouter()

	w := performRequest(router, "GET", "/")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("content-type"), "text/html")
}

func TestSitemap(t *testing.T) {
	router := SetupRouter()

	w := performRequest(router, "GET", "/sitemap.xml")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("content-type"), "text/xml")
}
