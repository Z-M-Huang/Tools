package main

import (
	"time"

	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
)

var appList []*webdata.AppCardList

func init() {
	getAnalyticTools()
}

func getAnalyticTools() {
	tools := &webdata.AppCardList{
		Category: "Analytic Tools",
	}

	kelly := &webdata.AppCard{
		FontsAwesomeTag: `<i class="fas fa-coins"></i>`,
		Link:            "/app/kelly-criterion",
		Title:           "Kelly Criterion",
		Description:     "Simulator for Kelly criterion. Kelly Criterion is a formula for sizing bets or investments from which the investor expects a positive return.",
		Usage:           0,
		Liked:           0,
	}

	tools.AppCards = append(tools.AppCards, kelly)
	appList = append(appList, tools)
}

func loadAppCardsUsage() {
	var tempAppList []*webdata.AppCardList
	copy(tempAppList, appList)
	for _, category := range tempAppList {
		for _, appCard := range category.AppCards {
			app := &dbentity.Application{}
			if db := utils.DB.Where(dbentity.Application{
				Name: appCard.Title,
			}).First(&app); db.RecordNotFound() {
				//Not found, let's insert
				app.Name = appCard.Title
				app.Usage = 0
				app.Liked = 0
				if dbIns := utils.DB.Save(app).Scan(&app); dbIns.Error != nil {
					utils.Logger.Sugar().Errorf("Failed to insert app %s into database. %s", appCard.Title, dbIns.Error.Error())
				}
			} else if app != nil {
				appCard.Usage = app.Usage
				appCard.Liked = app.Liked
			} else {
				utils.Logger.Sugar().Errorf("Failed to load app data from db for %s. %s", appCard.Title, db.Error.Error())
			}
		}
	}
	appList = tempAppList
}

func routineUpdateAppCardUsage(duration time.Duration) {
	for {
		loadAppCardsUsage()
		time.Sleep(duration)
	}
}
