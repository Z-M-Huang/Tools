package main

import (
	"compress/gzip"
	"context"
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
	"github.com/julienschmidt/httprouter"
)

func apiAuthHandler(requireToken bool, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		claim, err := getClaimFromHeaderAndRenew(w, r)
		if requireToken && (err != nil || claim == nil) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), utils.ClaimCtxKey, claim)
		ctx = context.WithValue(ctx, utils.ResponseCtxKey, &data.Response{})
		next(w, r.WithContext(ctx), ps)
	}
}

func pageHandler(requireToken bool, next httprouter.Handle) httprouter.Handle {
	return gzipHandler(pageAuthHandler(requireToken, pageStyleHandler(next)))
}

func apiHandler(requireToken bool, next httprouter.Handle) httprouter.Handle {
	return gzipHandler(apiAuthHandler(requireToken, next))
}

func pageStyleHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)
		style := ""
		cookie, err := r.Cookie(utils.PageStyleCookieKey)
		if err == nil {
			style = cookie.Value
		} else {
			logic.SetCookie(w, utils.PageStyleCookieKey, "default", time.Now().AddDate(100, 0, 0))
		}
		response.SetNavStyleName(logic.GetPageStyle(style))
		next(w, r, ps)
	}
}

func pageAuthHandler(requireToken bool, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		response := &data.Response{
			Header: &data.HeaderData{
				ResourceVersion: utils.Config.ResourceVersion,
			},
		}

		claim, err := getClaimFromCookieAndRenew(w, r)
		if requireToken && (err != nil || claim == nil) {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusTemporaryRedirect)
			return
		} else if claim != nil {
			response.SetLogin(&data.LoginData{
				Username: claim.Subject,
				ImageURL: claim.ImageURL,
			})
		}
		ctx := context.WithValue(r.Context(), utils.ClaimCtxKey, claim)
		ctx = context.WithValue(ctx, utils.ResponseCtxKey, response)
		next(w, r.WithContext(ctx), ps)
	}
}

func gzipHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r, ps)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := logic.GzipResponseWriter{Writer: gz, ResponseWriter: w}
		next(gzw, r, ps)
	}
}

func getClaimFromCookieAndRenew(w http.ResponseWriter, r *http.Request) (*data.JWTClaim, error) {
	cookie, err := r.Cookie(utils.SessionTokenKey)
	if err != nil {
		return nil, nil
	}
	claim, err := isTokenValid(cookie.Value)
	if err != nil {
		return nil, err
	}
	if time.Now().UTC().Sub(time.Unix(claim.ExpiresAt, 0)).Hours() < 24 {
		tokenStr, expiresAt, err := userlogic.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			logic.SetCookie(w, utils.SessionTokenKey, tokenStr, expiresAt)
		}
	}
	return claim, nil
}

func getClaimFromHeaderAndRenew(w http.ResponseWriter, r *http.Request) (*data.JWTClaim, error) {
	token := r.Header.Get("Authorization")
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
			logic.SetCookie(w, utils.SessionTokenKey, tokenStr, expiresAt)
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
	router := httprouter.New()

	assetsServer := http.FileServer(http.Dir("assets/"))
	vendorServer := http.FileServer(http.Dir("node_modules/"))

	router.GET(fmt.Sprintf("/assets/%s/*filepath", utils.Config.ResourceVersion),
		gzipHandler(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Set("Cache-Control", "public, max-age=604800")
			r.URL.Path = ps.ByName("filepath")
			assetsServer.ServeHTTP(w, r)
		}))

	router.GET(fmt.Sprintf("/vendor/%s/*filepath", utils.Config.ResourceVersion),
		gzipHandler(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Set("Cache-Control", "public, max-age=604800")
			r.URL.Path = ps.ByName("filepath")
			vendorServer.ServeHTTP(w, r)
		}))

	if utils.Config.SitemapConfig.GenerateSitemap {
		sm := utils.BuildSitemap()
		robotsContent := utils.GetRobotsTxt()
		err := utils.WriteContentToFile(robotsContent, "assets/robots.txt")
		if err != nil {
			utils.Logger.Fatal("Failed to create robots.txt")
		}

		router.GET("/sitemap.xml", gzipHandler(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Write(sm.XMLContent())
			return
		}))
		router.GET("/robots.txt", gzipHandler(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			http.ServeFile(w, r, "assets/robots.txt")
		}))

		if !utils.Config.IsDebug {
			go func() {
				time.Sleep(2000)
				utils.PingSearchEngines()
			}()
		}
	}

	router.GET("/", pageHandler(false, pages.HomePage))
	router.GET("/signup", pageHandler(false, pages.SignupPage))
	router.GET("/login", pageHandler(false, pages.LoginPage))
	router.GET("/account", pageHandler(true, pages.AccountPage))

	router.GET("/google_login", api.GoogleLogin)
	router.GET("/google_oauth", api.GoogleCallback)
	router.POST("/api/login", apiAuthHandler(false, api.Login))
	router.POST("/api/signup", apiAuthHandler(false, api.SignUp))
	router.POST("/api/account/update/password", apiAuthHandler(true, api.UpdatePassword))

	//app
	router.GET("/app/:name", pageHandler(false, pages.RenderApplicationPage))

	//app api
	router.POST("/api/kelly-criterion/simulate", apiHandler(false, appApis.KellyCriterionSimulate))
	router.POST("/api/hilo-simulator/simulate", apiHandler(false, appApis.HILOSimulate))
	router.POST("/api/hilo-simulator/verify", apiHandler(false, appApis.HILOVerify))
	router.POST("/api/dns-lookup/lookup", apiHandler(false, appApis.DNSLookup))
	router.POST("/api/string/encodedecode", apiHandler(false, appApis.EncodeDecode))
	router.POST("/app/:name/like", apiHandler(true, appApis.Like))
	router.POST("/app/:name/dislike", apiHandler(true, appApis.Dislike))

	utils.Logger.Fatal(http.ListenAndServe(":80", router).Error())
}
