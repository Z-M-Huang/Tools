package applicationlogic

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/copier"
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

	if usedAppCookie.Value == "" {
		return usedApps, nil
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

//GetApplicationWithLiked get application with liked populated
func GetApplicationWithLiked(user *dbentity.User) []*webdata.AppCategory {
	if user != nil && len(user.LikedApps) > 0 {
		appList := GetNewInstancesOfAppList()
		for _, category := range appList {
			for _, app := range category.AppCards {
				for _, likedApp := range user.LikedApps {
					if app.Title == likedApp {
						app.Liked = true
					}
				}
			}
		}
		return appList
	}
	return nil
}

//GetNewInstancesOfAppList get new instance
func GetNewInstancesOfAppList() []*webdata.AppCategory {
	var categories []*webdata.AppCategory
	for _, category := range utils.AppList {
		c := &webdata.AppCategory{
			Category: category.Category,
		}
		for _, app := range category.AppCards {
			temp := &webdata.AppCard{}
			copier.Copy(temp, app)
			c.AppCards = append(c.AppCards, temp)
		}

		categories = append(categories, c)
	}
	return categories
}
