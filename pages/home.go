package pages

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/julienschmidt/httprouter"
)

//HomePage home page /
func HomePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	if !(claim == nil) {
		user := &dbentity.User{
			Email: claim.Id,
		}
		err := userlogic.Find(utils.DB, user)
		if err == nil {
			if len(user.LikedApps) > 0 {
				response.Data = applicationlogic.GetApplicationWithLiked(user)
			} else {
				response.Data = utils.AppList
			}
		} else {
			response.Data = utils.AppList
		}
	} else {
		response.Data = utils.AppList
	}

	response.Header.Title = "Fun Apps"
	utils.Templates.ExecuteTemplate(w, "homepage.gohtml", response)
}

//GetPageStyle get page style
func GetPageStyle(style string) *data.PageStyleData {
	switch style {
	case "cerulean":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cerulean/bootstrap.min.css",
			Integrity: "sha384-LV/SIoc08vbV9CCeAwiz7RJZMI5YntsH8rGov0Y2nysmepqMWVvJqds6y0RaxIXT",
		}
	case "cosmo":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cosmo/bootstrap.min.css",
			Integrity: "sha384-qdQEsAI45WFCO5QwXBelBe1rR9Nwiss4rGEqiszC+9olH1ScrLrMQr1KmDR964uZ",
		}
	case "cyborg":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cyborg/bootstrap.min.css",
			Integrity: "sha384-l7xaoY0cJM4h9xh1RfazbgJVUZvdtyLWPueWNtLAphf/UbBgOVzqbOTogxPwYLHM",
		}
	case "darkly":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/darkly/bootstrap.min.css",
			Integrity: "sha384-rCA2D+D9QXuP2TomtQwd+uP50EHjpafN+wruul0sXZzX/Da7Txn4tB9aLMZV4DZm"}
	case "flatly":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/flatly/bootstrap.min.css",
			Integrity: "sha384-yrfSO0DBjS56u5M+SjWTyAHujrkiYVtRYh2dtB3yLQtUz3bodOeialO59u5lUCFF"}
	case "journal":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/journal/bootstrap.min.css",
			Integrity: "sha384-0d2eTc91EqtDkt3Y50+J9rW3tCXJWw6/lDgf1QNHLlV0fadTyvcA120WPsSOlwga"}
	case "litera":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/litera/bootstrap.min.css",
			Integrity: "sha384-pLgJ8jZ4aoPja/9zBSujjzs7QbkTKvKw1+zfKuumQF9U+TH3xv09UUsRI52fS+A6"}
	case "lumen":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/lumen/bootstrap.min.css",
			Integrity: "sha384-715VMUUaOfZR3/rWXZJ9agJ/udwSXGEigabzUbJm2NR4/v5wpDy8c14yedZN6NL7"}
	case "lux":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/lux/bootstrap.min.css",
			Integrity: "sha384-oOs/gFavzADqv3i5nCM+9CzXe3e5vXLXZ5LZ7PplpsWpTCufB7kqkTlC9FtZ5nJo"}
	case "materia":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/materia/bootstrap.min.css",
			Integrity: "sha384-1tymk6x9Y5K+OF0tlmG2fDRcn67QGzBkiM3IgtJ3VrtGrIi5ryhHjKjeeS60f1FA"}
	case "minty":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/minty/bootstrap.min.css",
			Integrity: "sha384-4HfFay3AYJnEmbgRzxYWJk/Ka5jIimhB/Fssk7NGT9Tj3rkEChpSxLK0btOGzf2I"}
	case "pulse":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/pulse/bootstrap.min.css",
			Integrity: "sha384-FnujoHKLiA0lyWE/5kNhcd8lfMILbUAZFAT89u11OhZI7Gt135tk3bGYVBC2xmJ5"}
	case "sandstone":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/sandstone/bootstrap.min.css",
			Integrity: "sha384-ABdnjefqVzESm+f9z9hcqx2cvwvDNjfrwfW5Le9138qHCMGlNmWawyn/tt4jR4ba"}
	case "simplex":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/simplex/bootstrap.min.css",
			Integrity: "sha384-cRAmF0wErT4D9dEBc36qB6pVu+KmLh516IoGWD/Gfm6FicBbyDuHgS4jmkQB8u1a"}
	case "sketchy":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/sketchy/bootstrap.min.css",
			Integrity: "sha384-2kOE+STGAkgemIkUbGtoZ8dJLqfvJ/xzRnimSkQN7viOfwRvWseF7lqcuNXmjwrL"}
	case "slate":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/slate/bootstrap.min.css",
			Integrity: "sha384-G9YbB4o4U6WS4wCthMOpAeweY4gQJyyx0P3nZbEBHyz+AtNoeasfRChmek1C2iqV"}
	case "solar":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/solar/bootstrap.min.css",
			Integrity: "sha384-Ya0fS7U2c07GI3XufLEwSQlqhpN0ri7w/ujyveyqoxTJ2rFHJcT6SUhwhL7GM4h9"}
	case "spacelab":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/spacelab/bootstrap.min.css",
			Integrity: "sha384-nl8CRcNYOGaXP68QRJeVTXCWAhwqBhRha0QbuubX1qDQpGT3GgclpvyzkR2JzyfZ"}
	case "superhero":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/superhero/bootstrap.min.css",
			Integrity: "sha384-R/oa7KS0iDoHwdh4Gyl3/fU7pgvSCt7oyuQ79pkw+e+bMWD9dzJJa+Zqd+XJS0AD"}
	case "united":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/united/bootstrap.min.css",
			Integrity: "sha384-bzjLLgZOhgXbSvSc5A9LWWo/mSIYf7U7nFbmYIB2Lgmuiw3vKGJuu+abKoaTx4W6"}
	case "yeti":
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/yeti/bootstrap.min.css",
			Integrity: "sha384-bWm7zrSUE5E+21rA9qdH5frkMpXvqjQm/WJw4L5PYQLNHrywI5zs5saEjIcCdGu1"}
	default:
		return &data.PageStyleData{
			Link:      "https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css",
			Integrity: "sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk"}
	}
}
