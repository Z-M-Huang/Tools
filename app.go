package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/pages"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pageData := &data.Response{}
	var cardList []*webdata.AppCardList
	for i := 0; i < 5; i++ {
		cardCategory := &webdata.AppCardList{
			Category: strconv.Itoa(i),
		}
		for j := 0; j < 10; j++ {
			cardCategory.AppCards = append(cardCategory.AppCards, &webdata.AppCard{
				Link:        "#",
				Title:       fmt.Sprintf("Card Title %d-%d", i, j),
				Description: "Some quick example text to build on the card title and make up the bulk of the card's content.",
				Up:          10000,
				Saved:       10000,
			})
		}
		cardList = append(cardList, cardCategory)
	}
	pageData.Data = cardList
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	if !claim.IsNil() {
		pageData.Login = data.LoginData{
			Username: claim.Subject,
			ImageURL: claim.ImageURL,
		}
	}
	utils.Templates.ExecuteTemplate(w, "homepage.gohtml", pageData)
}

func apiAuthHandler(requireClaim bool, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		claim, err := getClaimFromHeaderAndRenew(w, r)
		if (err != nil || claim.IsNil()) && requireClaim {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		ctx := context.WithValue(r.Context(), utils.ClaimCtxKey, claim)
		next(w, r.WithContext(ctx), ps)
	}
}

func pageAuthHandler(requireClaim bool, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		pageData := &data.Response{}
		claim, err := getClaimFromCookieAndRenew(w, r)
		if (err != nil || claim.IsNil()) && requireClaim {
			pageData.Alert.IsDanger = true
			pageData.Alert.Message = "Please login first"
			utils.Templates.ExecuteTemplate(w, "login.gohtml", pageData)
			return
		}
		ctx := context.WithValue(r.Context(), utils.ClaimCtxKey, claim)
		next(w, r.WithContext(ctx), ps)
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
		tokenStr, expiresAt, err := utils.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			setAuthCookie(w, tokenStr, expiresAt)
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
		tokenStr, expiresAt, err := utils.GenerateJWTToken(claim.Audience, claim.Id, claim.Subject, claim.ImageURL)
		if err != nil {
			utils.Logger.Sugar().Errorf("failed to generate jwt token %s", err.Error())
		} else {
			setAuthCookie(w, tokenStr, expiresAt)
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

func setAuthCookie(w http.ResponseWriter, tokenStr string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:       utils.SessionTokenKey,
		Value:      tokenStr,
		Path:       "/",
		Domain:     utils.Config.Host,
		Expires:    expiresAt,
		RawExpires: expiresAt.String(),
	})
}

func main() {

	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/vendor/*filepath", http.Dir("node_modules/"))

	router.GET("/", pageAuthHandler(false, homePage))
	router.GET("/login", pageAuthHandler(false, pages.LoginPage))
	router.POST("/login", pageAuthHandler(false, pages.Login))
	router.GET("/google_login", pages.GoogleLogin)
	router.GET("/google_oauth", pages.GoogleCallback)

	router.GET("/account", pageAuthHandler(true, pages.AccountPage))

	utils.Logger.Fatal(http.ListenAndServe(":80", router).Error())
}
