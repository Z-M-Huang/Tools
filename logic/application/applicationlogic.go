package applicationlogic

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"
)

//Find find db application
func Find(tx *gorm.DB, a *dbentity.Application) error {
	if db := tx.Where(*a).First(&a); db.Error != nil {
		return db.Error
	}
	return nil
}

//Save save application
func Save(tx *gorm.DB, a *dbentity.Application) error {
	if db := tx.Save(a).Scan(&a); db.Error != nil {
		return db.Error
	}
	return nil
}

//GetApplicationUsed saved in cookie
func GetApplicationUsed(r *http.Request) ([]string, error) {
	var usedApps []string
	usedAppCookie, err := r.Cookie(utils.UsedTokenKey)
	if err == http.ErrNoCookie {
		return usedApps, nil
	} else if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(usedAppCookie.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(decoded, &usedApps)
	if err != nil {
		return nil, err
	}
	return usedApps, nil
}

//GetApplicationsByName get application by name
func GetApplicationsByName(name string) *webdata.AppCard {
	for _, category := range utils.AppList {
		for _, app := range category.AppCards {
			if app.Name == name {
				return app
			}
		}
	}
	return nil
}
