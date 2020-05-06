package api

import (
	"fmt"
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

//UpdatePassword api
func UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := &data.Response{}

	token := r.Header.Get("Authorization")
	if token == "" {
		resp.Alert.IsDanger = true
		resp.Alert.Message = "Unauthorized"
		resp.Data = false
		writeResponse(w, resp)
		return
	}

	claim := r.Context().Value(utils.ClaimCtxKey).(*data.JWTClaim)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse request data", http.StatusBadRequest)
		return
	}

	password := r.FormValue("newPassword")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		http.Error(w, "Password doesn't match.", http.StatusBadRequest)
		return
	} else if len(password) < 12 {
		http.Error(w, "Password has minimum length of 12", http.StatusBadRequest)
		return
	}

	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		dbUser := &dbentity.User{
			Email: claim.Id,
		}
		if db := tx.Where(dbUser).First(&dbUser); db.RecordNotFound() {
			return fmt.Errorf("user not found for email: %s", claim.Id)
		} else if db.Error != nil {
			return fmt.Errorf(fmt.Sprintf("failed to user in UpdatePassword %s", db.Error.Error()))
		}

		dbUser.Password = utils.HashAndSalt([]byte(password))

		if db := tx.Save(dbUser).Scan(&dbUser); db.Error != nil {
			return fmt.Errorf(fmt.Sprintf("failed to update user with new password %s", db.Error.Error()))
		}

		return nil
	})

	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, "Failed to update user information, please try again later.", http.StatusInternalServerError)
		return
	}
}
