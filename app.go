package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/api"
	appApis "github.com/Z-M-Huang/Tools/api/app"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/logic"
	"github.com/Z-M-Huang/Tools/pages"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func apiAuthHandler(requireToken bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := logic.GetClaimFromCookieAndRenew(c)
		if requireToken && (err != nil || claim == nil) {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Set(utils.ClaimCtxKey, claim)
		c.Set(utils.ResponseCtxKey, &data.Response{})
		c.Next()
	}
}

func pageStyleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := c.Keys[utils.ResponseCtxKey].(*data.Response)
		style := ""
		val, err := c.Cookie(utils.PageStyleCookieKey)
		if err == nil && val != "" {
			style = val
		} else {
			logic.SetCookie(c, utils.PageStyleCookieKey, "default", time.Now().AddDate(100, 0, 0), false)
		}
		response.SetNavStyleName(getPageStyle(style))
		c.Next()
	}
}

func pageAuthHandler(requireToken bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := &data.Response{
			Header: &data.HeaderData{
				ResourceVersion: data.Config.ResourceVersion,
			},
		}

		claim, err := logic.GetClaimFromCookieAndRenew(c)
		if requireToken && (err != nil || claim == nil) {
			c.Redirect(http.StatusTemporaryRedirect, "/login?redirect="+c.Request.URL.Path)
			c.Abort()
			return
		} else if claim != nil {
			response.SetLogin(&data.LoginData{
				Username: claim.Subject,
				ImageURL: claim.ImageURL,
			})
		}
		c.Set(utils.ClaimCtxKey, claim)
		c.Set(utils.ResponseCtxKey, response)
		c.Next()
	}
}

func getPageStyle(style string) *data.PageStyleData {
	switch style {
	case "cerulean":
		return &data.PageStyleData{
			Name:      "Cerulean",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cerulean/bootstrap.min.css",
			Integrity: "sha384-LV/SIoc08vbV9CCeAwiz7RJZMI5YntsH8rGov0Y2nysmepqMWVvJqds6y0RaxIXT",
		}
	case "cosmo":
		return &data.PageStyleData{
			Name:      "Cosmo",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cosmo/bootstrap.min.css",
			Integrity: "sha384-qdQEsAI45WFCO5QwXBelBe1rR9Nwiss4rGEqiszC+9olH1ScrLrMQr1KmDR964uZ",
		}
	case "cyborg":
		return &data.PageStyleData{
			Name:      "Cyborg",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/cyborg/bootstrap.min.css",
			Integrity: "sha384-l7xaoY0cJM4h9xh1RfazbgJVUZvdtyLWPueWNtLAphf/UbBgOVzqbOTogxPwYLHM",
		}
	case "darkly":
		return &data.PageStyleData{
			Name:      "Darkly",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/darkly/bootstrap.min.css",
			Integrity: "sha384-rCA2D+D9QXuP2TomtQwd+uP50EHjpafN+wruul0sXZzX/Da7Txn4tB9aLMZV4DZm"}
	case "flatly":
		return &data.PageStyleData{
			Name:      "Flatly",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/flatly/bootstrap.min.css",
			Integrity: "sha384-yrfSO0DBjS56u5M+SjWTyAHujrkiYVtRYh2dtB3yLQtUz3bodOeialO59u5lUCFF"}
	case "journal":
		return &data.PageStyleData{
			Name:      "Journal",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/journal/bootstrap.min.css",
			Integrity: "sha384-0d2eTc91EqtDkt3Y50+J9rW3tCXJWw6/lDgf1QNHLlV0fadTyvcA120WPsSOlwga"}
	case "litera":
		return &data.PageStyleData{
			Name:      "Litera",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/litera/bootstrap.min.css",
			Integrity: "sha384-pLgJ8jZ4aoPja/9zBSujjzs7QbkTKvKw1+zfKuumQF9U+TH3xv09UUsRI52fS+A6"}
	case "lumen":
		return &data.PageStyleData{
			Name:      "Lumen",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/lumen/bootstrap.min.css",
			Integrity: "sha384-715VMUUaOfZR3/rWXZJ9agJ/udwSXGEigabzUbJm2NR4/v5wpDy8c14yedZN6NL7"}
	case "lux":
		return &data.PageStyleData{
			Name:      "Lux",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/lux/bootstrap.min.css",
			Integrity: "sha384-oOs/gFavzADqv3i5nCM+9CzXe3e5vXLXZ5LZ7PplpsWpTCufB7kqkTlC9FtZ5nJo"}
	case "materia":
		return &data.PageStyleData{
			Name:      "Materia",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/materia/bootstrap.min.css",
			Integrity: "sha384-1tymk6x9Y5K+OF0tlmG2fDRcn67QGzBkiM3IgtJ3VrtGrIi5ryhHjKjeeS60f1FA"}
	case "minty":
		return &data.PageStyleData{
			Name:      "Minty",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/minty/bootstrap.min.css",
			Integrity: "sha384-4HfFay3AYJnEmbgRzxYWJk/Ka5jIimhB/Fssk7NGT9Tj3rkEChpSxLK0btOGzf2I"}
	case "pulse":
		return &data.PageStyleData{
			Name:      "Pulse",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/pulse/bootstrap.min.css",
			Integrity: "sha384-FnujoHKLiA0lyWE/5kNhcd8lfMILbUAZFAT89u11OhZI7Gt135tk3bGYVBC2xmJ5"}
	case "sandstone":
		return &data.PageStyleData{
			Name:      "Sandstone",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/sandstone/bootstrap.min.css",
			Integrity: "sha384-ABdnjefqVzESm+f9z9hcqx2cvwvDNjfrwfW5Le9138qHCMGlNmWawyn/tt4jR4ba"}
	case "simplex":
		return &data.PageStyleData{
			Name:      "Simplex",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/simplex/bootstrap.min.css",
			Integrity: "sha384-cRAmF0wErT4D9dEBc36qB6pVu+KmLh516IoGWD/Gfm6FicBbyDuHgS4jmkQB8u1a"}
	case "sketchy":
		return &data.PageStyleData{
			Name:      "Sketchy",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/sketchy/bootstrap.min.css",
			Integrity: "sha384-2kOE+STGAkgemIkUbGtoZ8dJLqfvJ/xzRnimSkQN7viOfwRvWseF7lqcuNXmjwrL"}
	case "slate":
		return &data.PageStyleData{
			Name:      "Salte",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/slate/bootstrap.min.css",
			Integrity: "sha384-G9YbB4o4U6WS4wCthMOpAeweY4gQJyyx0P3nZbEBHyz+AtNoeasfRChmek1C2iqV"}
	case "solar":
		return &data.PageStyleData{
			Name:      "Solar",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/solar/bootstrap.min.css",
			Integrity: "sha384-Ya0fS7U2c07GI3XufLEwSQlqhpN0ri7w/ujyveyqoxTJ2rFHJcT6SUhwhL7GM4h9"}
	case "spacelab":
		return &data.PageStyleData{
			Name:      "Spacelab",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/spacelab/bootstrap.min.css",
			Integrity: "sha384-nl8CRcNYOGaXP68QRJeVTXCWAhwqBhRha0QbuubX1qDQpGT3GgclpvyzkR2JzyfZ"}
	case "superhero":
		return &data.PageStyleData{
			Name:      "Superhero",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/superhero/bootstrap.min.css",
			Integrity: "sha384-R/oa7KS0iDoHwdh4Gyl3/fU7pgvSCt7oyuQ79pkw+e+bMWD9dzJJa+Zqd+XJS0AD"}
	case "united":
		return &data.PageStyleData{
			Name:      "United",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/united/bootstrap.min.css",
			Integrity: "sha384-bzjLLgZOhgXbSvSc5A9LWWo/mSIYf7U7nFbmYIB2Lgmuiw3vKGJuu+abKoaTx4W6"}
	case "yeti":
		return &data.PageStyleData{
			Name:      "Yeti",
			Link:      "https://stackpath.bootstrapcdn.com/bootswatch/4.4.1/yeti/bootstrap.min.css",
			Integrity: "sha384-bWm7zrSUE5E+21rA9qdH5frkMpXvqjQm/WJw4L5PYQLNHrywI5zs5saEjIcCdGu1"}
	default:
		return &data.PageStyleData{
			Name:      "Default",
			Link:      "https://stackpath.bootstrapcdn.com/bootstrap/4.5.0/css/bootstrap.min.css",
			Integrity: "sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk"}
	}
}

//SetupRouter setup gin router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetHTMLTemplate(utils.Templates)

	router.Static(fmt.Sprintf("/assets/%s", data.Config.ResourceVersion), "assets/")
	router.Static(fmt.Sprintf("/vendor/%s", data.Config.ResourceVersion), "node_modules/")
	router.StaticFile("/favicon.ico", "assets/img/favicon.ico")

	if data.Config.EnableSitemap {
		sm := BuildSitemap()
		robotsContent := GetRobotsTxt()
		err := utils.WriteContentToFile(robotsContent, "assets/robots.txt")
		if err != nil {
			utils.Logger.Fatal("Failed to create robots.txt")
		}
		router.StaticFile("robots.txt", "assets/robots.txt")

		router.GET("/sitemap.xml", func(c *gin.Context) {
			c.Writer.Write(sm.XMLContent())
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Header().Add("Content-Type", "text/xml; charset=UTF-8")
		})
		if !data.Config.IsDebug {
			go func() {
				time.Sleep(2000)
				PingSearchEngines()
			}()
		}
	}

	pageNoAuth := router.Group("/", pageAuthHandler(false), pageStyleHandler())
	pageAuthRequired := router.Group("/", pageAuthHandler(true), pageStyleHandler())
	apiNoAuth := router.Group("/api", apiAuthHandler(false))
	apiAuthRequired := router.Group("/api", apiAuthHandler(true))

	pageNoAuth.GET("/", pages.HomePage)
	pageNoAuth.GET("/signup", pages.SignupPage)
	pageNoAuth.GET("/login", pages.LoginPage)
	pageAuthRequired.GET("/account", pages.AccountPage)

	router.GET("/google_login", api.GoogleLogin)
	router.GET("/google_oauth", api.GoogleCallback)
	apiNoAuth.POST("/login", api.Login)
	apiNoAuth.POST("/logout", api.Logout)
	apiNoAuth.POST("/signup", api.SignUp)
	apiAuthRequired.POST("/account/update/password", api.UpdatePassword)

	//app
	pageNoAuth.GET("/app/:name", pages.RenderApplicationPage)
	pageNoAuth.GET("/app/:name/:id", pages.RenderApplicationPage)

	//app api
	router.Any("/api/request-bin/receive/:id", appApis.RequestIn)
	apiNoAuth.POST("/kelly-criterion/simulate", appApis.KellyCriterionSimulate)
	apiNoAuth.POST("/hilo-simulator/simulate", appApis.HILOSimulate)
	apiNoAuth.POST("/hilo-simulator/verify", appApis.HILOVerify)
	apiNoAuth.POST("/dns-lookup/lookup", appApis.DNSLookup)
	apiNoAuth.POST("/string/encodedecode", appApis.EncodeDecode)
	apiNoAuth.POST("/request-bin/create", appApis.CreateRequestBin)
	apiNoAuth.POST("/qr-code/create", appApis.CreateQRCode)
	apiAuthRequired.POST("/app/:name/like", appApis.Like)
	apiAuthRequired.POST("/app/:name/dislike", appApis.Dislike)

	return router
}

func main() {
	utils.Logger.Fatal(http.ListenAndServe(":80", SetupRouter()).Error())
}
