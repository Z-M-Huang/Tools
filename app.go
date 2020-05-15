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
		claim, err := logic.GetClaimFromHeaderAndRenew(c)
		if requireToken && (err != nil || claim == nil) {
			c.String(http.StatusUnauthorized, "Unauthorized")
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
			logic.SetCookie(c, utils.PageStyleCookieKey, "default", time.Now().AddDate(100, 0, 0), true)
		}
		response.SetNavStyleName(logic.GetPageStyle(style))
		c.Next()
	}
}

func pageAuthHandler(requireToken bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := &data.Response{
			Header: &data.HeaderData{
				ResourceVersion: utils.Config.ResourceVersion,
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

func main() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetHTMLTemplate(utils.Templates)

	router.Static(fmt.Sprintf("/assets/%s", utils.Config.ResourceVersion), "assets/")
	router.Static(fmt.Sprintf("/vendor/%s", utils.Config.ResourceVersion), "node_modules/")
	router.StaticFile("/favicon.ico", "assets/img/favicon.ico")

	if utils.Config.EnableSitemap {
		sm := utils.BuildSitemap()
		robotsContent := utils.GetRobotsTxt()
		err := utils.WriteContentToFile(robotsContent, "assets/robots.txt")
		if err != nil {
			utils.Logger.Fatal("Failed to create robots.txt")
		}
		router.StaticFile("robots.txt", "assets/robots.txt")

		router.GET("/sitemap.xml", func(c *gin.Context) {
			c.Writer.Write(sm.XMLContent())
		})
		if !utils.Config.IsDebug {
			go func() {
				time.Sleep(2000)
				utils.PingSearchEngines()
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
	apiNoAuth.POST("/kelly-criterion/simulate", appApis.KellyCriterionSimulate)
	apiNoAuth.POST("/hilo-simulator/simulate", appApis.HILOSimulate)
	apiNoAuth.POST("/hilo-simulator/verify", appApis.HILOVerify)
	apiNoAuth.POST("/dns-lookup/lookup", appApis.DNSLookup)
	apiNoAuth.POST("/string/encodedecode", appApis.EncodeDecode)
	apiNoAuth.POST("/request-bin/create", appApis.CreateRequestBin)
	apiAuthRequired.POST("/app/:name/like", appApis.Like)
	apiAuthRequired.POST("/app/:name/dislike", appApis.Dislike)

	utils.Logger.Fatal(http.ListenAndServe(":80", router).Error())
}
