package app

import (
	"fmt"
	"net/http"

	"github.com/Z-M-Huang/Tools/api"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	userlogic "github.com/Z-M-Huang/Tools/logic/user"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

//Like /app/:name/like
func Like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Only logged in user can access this
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)

	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	appCard := applicationlogic.GetApplicationsByName(name)
	if appCard == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	user := &dbentity.User{
		Email: claim.Id,
	}
	err := userlogic.Find(utils.DB, user)
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("oh boy... There is a user doesn't found in database but have a token. Email: %s", claim.Id)
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		api.WriteResponse(w, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		api.WriteResponse(w, response)
		return
	}

	found := false
	for _, likedApp := range user.LikedApps {
		if likedApp == appCard.Title {
			found = true
			break
		}
	}

	if !found {
		user.LikedApps = append(user.LikedApps, appCard.Title)
		err = utils.DB.Transaction(func(tx *gorm.DB) error {
			err = userlogic.Save(utils.DB, user)
			if err != nil {
				return fmt.Errorf("failed to add liked app %s to user %s", appCard.Title, claim.Id)
			}
			err = addApplicationLike(tx, appCard)
			if err != nil {
				return fmt.Errorf("failed to add like to app %s", appCard.Title)
			}
			return nil
		})
		if err != nil {
			api.WriteUnexpectedError(w, response)
			utils.Logger.Error(err.Error())
			return
		}
	}
	response.SetAlert(&data.AlertData{
		IsSuccess: true,
		Message:   "Application saved! Thank you for support.",
	})
	response.Data = appCard.AmountLiked
	api.WriteResponse(w, response)
}

//Dislike /app/:name/dislike
func Dislike(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Only logged in user can access this
	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)
	response := r.Context().Value(utils.ResponseCtxKey).(*data.Response)

	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	appCard := applicationlogic.GetApplicationsByName(name)
	if appCard == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	user := &dbentity.User{
		Email: claim.Id,
	}
	err := userlogic.Find(utils.DB, user)
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("oh boy... There is a user doesn't found in database but have a token. Email: %s", claim.Id)
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		api.WriteResponse(w, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		api.WriteResponse(w, response)
		return
	}

	index := -1
	for i, likedApp := range user.LikedApps {
		if likedApp == appCard.Title {
			index = i
			break
		}
	}

	if index > -1 {
		user.LikedApps[index] = user.LikedApps[len(user.LikedApps)-1]
		user.LikedApps = user.LikedApps[:len(user.LikedApps)-1]
		err = utils.DB.Transaction(func(tx *gorm.DB) error {
			err = userlogic.Save(utils.DB, user)
			if err != nil {
				return fmt.Errorf("failed to remove liked app %s from user %s", appCard.Title, claim.Id)
			}
			err = removeApplicationLike(tx, appCard)
			if err != nil {
				return fmt.Errorf("failed to remove like from app %s", appCard.Title)
			}
			return nil
		})
		if err != nil {
			api.WriteUnexpectedError(w, response)
			utils.Logger.Error(err.Error())
			return
		}
	}
	response.SetAlert(&data.AlertData{
		IsInfo:  true,
		Message: "If there are anything you want us to improve about this app. Please let us know on the github bug tracker.",
	})
	response.Data = appCard.AmountLiked
	api.WriteResponse(w, response)
}

func addApplicationLike(tx *gorm.DB, app *webdata.AppCard) error {
	dbApp := &dbentity.Application{
		Name: app.Title,
	}
	err := applicationlogic.Find(tx, dbApp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	dbApp.Liked++
	err = applicationlogic.Save(tx, dbApp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	app.AmountLiked = dbApp.Liked
	return nil
}

func removeApplicationLike(tx *gorm.DB, app *webdata.AppCard) error {
	dbApp := &dbentity.Application{
		Name: app.Title,
	}
	err := applicationlogic.Find(tx, dbApp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	dbApp.Liked--
	err = applicationlogic.Save(tx, dbApp)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	app.AmountLiked = dbApp.Liked
	return nil
}
