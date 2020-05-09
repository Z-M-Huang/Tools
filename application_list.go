package main

import (
	"sort"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	"github.com/Z-M-Huang/Tools/utils"
)

func init() {
	getAnalyticTools()
	loadAppCardsUsage()
}

func getAnalyticTools() {
	tools := &webdata.AppCardList{
		Category: "Analytic Tools",
	}

	kelly := &webdata.AppCard{
		Name:            "kelly-criterion",
		TemplateName:    "kelly_criterion.gohtml",
		FontsAwesomeTag: `fas fa-coins`,
		Link:            "/app/kelly-criterion",
		Title:           "Kelly Criterion",
		Description:     "Simulator for Kelly criterion. Kelly Criterion is a formula for sizing bets or investments from which the investor expects a positive return.",
		Usage:           0,
		Liked:           0,
	}
	tools.AppCards = append(tools.AppCards, kelly)

	martingale := &webdata.AppCard{
		Name:            "martingale",
		TemplateName:    "martingale.gohtml",
		FontsAwesomeTag: `fas fa-comments-dollar`,
		Link:            "/app/martingale",
		Title:           "Martingale",
		Description:     "Simulator for Martingale. The simplest of these strategies was designed for a game in which the gambler wins the stake if a coin comes up heads and loses it if the coin comes up tails.",
		Usage:           0,
		Liked:           0,
	}
	tools.AppCards = append(tools.AppCards, martingale)

	sortAppCardSlice(tools.AppCards)
	utils.AppList = append(utils.AppList, tools)
}

func loadAppCardsUsage() {
	for _, category := range utils.AppList {
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
}

func routineUpdateAppCardUsage(duration time.Duration) {
	for {
		loadAppCardsUsage()
		time.Sleep(duration)
	}
}

func sortAppCardSlice(appCards []*webdata.AppCard) {
	sort.Slice(appCards, func(i, j int) bool {
		return strings.ToLower(appCards[i].Title) < strings.ToLower(appCards[j].Title)
	})
}
