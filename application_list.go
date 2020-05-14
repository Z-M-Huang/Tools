package main

import (
	"sort"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data/dbentity"
	"github.com/Z-M-Huang/Tools/data/webdata"
	applicationlogic "github.com/Z-M-Huang/Tools/logic/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"
)

func init() {
	getAnalyticTools()
	getConverterTools()
	getLookupTools()
	loadAppCardsUsage()
}

func getAnalyticTools() {
	tools := &webdata.AppCategory{
		Category: "Analytic Tools",
	}

	kelly := newAppCart("kelly-criterion", "kelly_criterion.gohtml", "", "fas fa-coins",
		"/app/kelly-criterion", "Kelly Criterion", "Simulator for Kelly criterion. Kelly Criterion is a formula for sizing bets or investments from which the investor expects a positive return.")
	tools.AppCards = append(tools.AppCards, kelly)

	betSimulator := newAppCart("hilo-simulator", "hilo_simulator.gohtml", "", "fas fa-coins",
		"/app/hilo-simulator", "HiLo Simulator", "Simulate online hi/low betting website result(Provably fair only).")

	tools.AppCards = append(tools.AppCards, betSimulator)

	sortAppCardSlice(tools.AppCards)
	utils.AppList = append(utils.AppList, tools)
}

func getConverterTools() {
	tools := &webdata.AppCategory{
		Category: "Converter",
	}

	encoderDecoder := newAppCart("string-encoder-decoder", "string_encoder_decoder.gohtml", "", "fas fa-receipt",
		"/app/string-encoder-decoder", "Encoder Decoder", "Convert string encoding based on request.")
	tools.AppCards = append(tools.AppCards, encoderDecoder)

	sortAppCardSlice(tools.AppCards)
	utils.AppList = append(utils.AppList, tools)
}

func getLookupTools() {
	tools := &webdata.AppCategory{
		Category: "Lookup",
	}

	dnsLookup := newAppCart("dns-lookup", "dns_lookup.gohtml", "", "fas fa-receipt",
		"/app/dns-lookup", "DNS Lookup", "Lookup given domain's DNS record (A, CNAME, PTR, NS, MX, TXT, and etc.")
	tools.AppCards = append(tools.AppCards, dnsLookup)

	sortAppCardSlice(tools.AppCards)
	utils.AppList = append(utils.AppList, tools)
}

func loadAppCardsUsage() {
	for _, category := range utils.AppList {
		for _, appCard := range category.AppCards {
			app := &dbentity.Application{
				Name: appCard.Title,
			}
			err := applicationlogic.Find(utils.DB, app)
			if err == gorm.ErrRecordNotFound {
				app.Name = appCard.Title
				app.Usage = 0
				app.Liked = 0
				err = applicationlogic.Save(utils.DB, app)
				if err != nil {
					utils.Logger.Sugar().Errorf("Failed to insert app %s into database. %s", appCard.Title, err.Error())
				}
			} else if err != nil {
				utils.Logger.Sugar().Errorf("Failed to load app data from db for %s. %s", appCard.Title, err.Error())
			} else {
				appCard.AmountUsed = app.Usage
				appCard.AmountLiked = app.Liked
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
