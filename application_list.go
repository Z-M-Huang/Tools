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

	kelly := newAppCart("kelly-criterion", "kelly_criterion.gohtml", "", "fas fa-coins",
		"/app/kelly-criterion", "Kelly Criterion", "Simulator for Kelly criterion. Kelly Criterion is a formula for sizing bets or investments from which the investor expects a positive return.")
	tools.AppCards = append(tools.AppCards, kelly)

	betSimulator := newAppCart("bet-simulator", "bet_simulator.gohtml", "", "fas fa-coins",
		"/app/bet-simulator", "Bet Simulator", "Simulate online betting website result(Provably fair only).")

	tools.AppCards = append(tools.AppCards, betSimulator)

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
				appCard.AmountUsed = app.Usage
				appCard.AmountLiked = app.Liked
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

func newAppCart(name, templateName, imageURL, fontsAwesomeTag, link, title, description string) *webdata.AppCard {
	return &webdata.AppCard{
		Name:            name,
		TemplateName:    templateName,
		FontsAwesomeTag: fontsAwesomeTag,
		Link:            link,
		Title:           title,
		Description:     description,
		Used:            false,
		AmountUsed:      0,
		Liked:           false,
		AmountLiked:     0,
	}
}
