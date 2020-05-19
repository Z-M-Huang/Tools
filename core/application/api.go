package application

import (
	"fmt"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//API api endpoint
type API struct{}

//Like /app/:name/like
func (API) Like(c *gin.Context) {
	//Only logged in user can access this
	claim := account.GetClaimInContext(c.Keys)
	response := core.GetResponseInContext(c.Keys)

	name := c.Param("name")
	if name == "" {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	appCard := GetApplicationsByName(name)
	if appCard == nil {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	user := &db.User{
		Email: claim.Id,
	}
	err := user.Find()
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("oh boy... There is a user doesn't found in database but have a token. Email: %s", claim.Id)
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
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
		err = db.DoTransaction(func(tx *gorm.DB) error {
			err = user.SaveWithTx(tx)
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
			core.WriteUnexpectedError(c, response)
			utils.Logger.Error(err.Error())
			return
		}
		go ReloadAppList()
	}
	response.SetAlert(&data.AlertData{
		IsSuccess: true,
		Message:   "Application saved! Thank you for support.",
	})
	response.Data = appCard.AmountLiked
	core.WriteResponse(c, 200, response)
}

//Dislike /app/:name/dislike
func (API) Dislike(c *gin.Context) {
	//Only logged in user can access this
	claim := c.Keys[utils.ClaimCtxKey].(*account.JWTClaim)
	response := c.Keys[utils.ResponseCtxKey].(*data.Response)

	name := c.Param("name")
	if name == "" {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	appCard := GetApplicationsByName(name)
	if appCard == nil {
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	user := &db.User{
		Email: claim.Id,
	}
	err := user.Find()
	if err == gorm.ErrRecordNotFound {
		utils.Logger.Sugar().Errorf("oh boy... There is a user doesn't found in database but have a token. Email: %s", claim.Id)
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&data.AlertData{
			IsDanger: true,
			Message:  "User not found",
		})
		core.WriteResponse(c, 400, response)
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
		err = db.DoTransaction(func(tx *gorm.DB) error {
			err = user.Save()
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
			core.WriteUnexpectedError(c, response)
			utils.Logger.Error(err.Error())
			return
		}
		go ReloadAppList()
	}
	response.SetAlert(&data.AlertData{
		IsInfo:  true,
		Message: "If there are anything you want us to improve about this app. Please let us know on the github bug tracker.",
	})
	response.Data = appCard.AmountLiked
	core.WriteResponse(c, 200, response)
}

func addApplicationLike(tx *gorm.DB, app *AppCard) error {
	dbApp := &db.Application{
		Name: app.Title,
	}
	err := dbApp.FindWithTx(tx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	dbApp.Liked++
	err = dbApp.SaveWithTx(tx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	app.AmountLiked = dbApp.Liked
	return nil
}

func removeApplicationLike(tx *gorm.DB, app *AppCard) error {
	dbApp := &db.Application{
		Name: app.Title,
	}
	err := dbApp.FindWithTx(tx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	dbApp.Liked--
	err = dbApp.SaveWithTx(tx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	app.AmountLiked = dbApp.Liked
	return nil
}
