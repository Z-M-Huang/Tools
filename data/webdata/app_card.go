package webdata

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Z-M-Huang/Tools/data/constval"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/jinzhu/copier"
)

//AppList in home page
var AppList []*AppCategory

//AppCategory card list for single category
type AppCategory struct {
	Category string
	AppCards []*AppCard
}

//AppCard used to render app card in pages
type AppCard struct {
	Name            string
	TemplateName    string
	ImageURL        string
	FontsAwesomeTag string
	Link            string
	Title           string
	Description     string
	Used            bool
	AmountUsed      uint64
	Liked           bool
	AmountLiked     uint64
}

//GetApplicationUsed saved in cookie
func GetApplicationUsed(r *http.Request) ([]string, error) {
	var usedApps []string
	usedAppCookie, err := r.Cookie(constval.UsedTokenKey)
	if err == http.ErrNoCookie {
		return usedApps, nil
	} else if err != nil {
		return nil, err
	}

	if usedAppCookie.Value == "" {
		return usedApps, nil
	}

	queryDecoded, err := url.QueryUnescape(usedAppCookie.Value)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(queryDecoded)
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
func GetApplicationsByName(name string) *AppCard {
	for _, category := range AppList {
		for _, app := range category.AppCards {
			if app.Name == name {
				return app
			}
		}
	}
	return nil
}

//GetApplicationWithLiked get application with liked populated
func GetApplicationWithLiked(user *db.User) []*AppCategory {
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
func GetNewInstancesOfAppList() []*AppCategory {
	var categories []*AppCategory
	for _, category := range AppList {
		c := &AppCategory{
			Category: category.Category,
		}
		for _, app := range category.AppCards {
			temp := &AppCard{}
			copier.Copy(temp, app)
			c.AppCards = append(c.AppCards, temp)
		}

		categories = append(categories, c)
	}
	return categories
}
