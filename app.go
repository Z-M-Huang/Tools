package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/api"
	appApis "github.com/Z-M-Huang/Tools/api/app"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/logic"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/pages"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func apiAuthHandler(requireToken bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := getClaimFromHeaderAndRenew(c)
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
			logic.SetCookie(c, utils.PageStyleCookieKey, "default", time.Now().AddDate(100, 0, 0))
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

		claim, err := getClaimFromCookieAndRenew(c)
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

func getClaimFromCookieAndRenew(c *gin.Context) (*data.JWTClaim, error) {
	val, err := c.Cookie(utils.SessionTokenKey)
	if err != nil || val == "" {
		return nil, nil
	}
	claim, err := isTokenValid(val)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().Sub(time.Unix(claim.ExpiresAt, 0)).Hours() < 24 {
		tokenStr, expiresAt, err := userlogic.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			logic.SetCookie(c, utils.SessionTokenKey, tokenStr, expiresAt)
		}
	}
	return claim, nil
}

func getClaimFromHeaderAndRenew(c *gin.Context) (*data.JWTClaim, error) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.Contains(token, "Bearer ") {
		return nil, errors.New("Unauthorized")
	}

	token = strings.ReplaceAll(token, "Bearer ", "")
	claim, err := isTokenValid(token)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}
	if time.Now().UTC().Sub(time.Unix(claim.ExpiresAt, 0)).Hours() < 24 {
		tokenStr, expiresAt, err := userlogic.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			logic.SetCookie(c, utils.SessionTokenKey, tokenStr, expiresAt)
		}
	}
	return claim, nil
}

func isTokenValid(token string) (*data.JWTClaim, error) {
	claims := &data.JWTClaim{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.Config.JwtKey, nil
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("Unauthenticated")
	}

	if !tkn.Valid || !claims.VerifyIssuer(utils.Config.Host, true) {
		return nil, fmt.Errorf("Invalid Token")
	}

	return claims, nil
}

func main() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static(fmt.Sprintf("/assets/%s", utils.Config.ResourceVersion), "assets/")
	router.Static(fmt.Sprintf("/vendor/%s", utils.Config.ResourceVersion), "node_modules/")

	if utils.Config.SitemapConfig.GenerateSitemap {
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
	apiNoAuth.POST("/signup", api.SignUp)
	apiAuthRequired.POST("/account/update/password", api.UpdatePassword)

	//app
	pageNoAuth.GET("/app/:name", pages.RenderApplicationPage)

	//app api
	apiNoAuth.POST("/kelly-criterion/simulate", appApis.KellyCriterionSimulate)
	apiNoAuth.POST("/hilo-simulator/simulate", appApis.HILOSimulate)
	apiNoAuth.POST("/hilo-simulator/verify", appApis.HILOVerify)
	apiNoAuth.POST("/dns-lookup/lookup", appApis.DNSLookup)
	apiNoAuth.POST("/string/encodedecode", appApis.EncodeDecode)
	apiAuthRequired.POST("/app/:name/like", appApis.Like)
	apiAuthRequired.POST("/app/:name/dislike", appApis.Dislike)

	utils.Logger.Fatal(http.ListenAndServe(":80", router).Error())
}
