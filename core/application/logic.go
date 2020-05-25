package application

import (
	"sort"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/jinzhu/gorm"
)

//GetApplicationsByName get application by name
func GetApplicationsByName(name string) *AppCard {
	for _, category := range GetAppList() {
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
		appList := GetAppList()
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

//GetAppList get app list from redis
func GetAppList() []*AppCategory {
	var categories []*AppCategory
	err := db.RedisGet(utils.RedisAppListKey, &categories)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	if categories == nil || len(categories) == 0 {
		var appList []*AppCategory
		appList = append(appList, getAnalyticTools(), getFormatTools(), getGeneratorTools(),
			getLookupTools(), getWebUtils())
		loadAppCardsUsage(appList)
		appList = append([]*AppCategory{getPopular(appList)}, appList...)
		err = db.RedisSetBytes(utils.RedisAppListKey, appList, 24*time.Hour)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
		categories = appList
	}
	return categories
}

//ReloadAppList reload app list
func ReloadAppList() {
	utils.Logger.Info("Reload AppList...")
	var categories []*AppCategory
	err := db.RedisDelete(utils.RedisAppListKey)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	var appList []*AppCategory
	appList = append(appList, getAnalyticTools(), getFormatTools(), getGeneratorTools(),
		getLookupTools(), getWebUtils())
	loadAppCardsUsage(appList)
	appList = append([]*AppCategory{getPopular(appList)}, appList...)
	categories = appList
	err = db.RedisSetBytes(utils.RedisAppListKey, categories, 24*time.Hour)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}

//AddApplicationUsage add usage
func AddApplicationUsage(app *AppCard) {
	app.AmountUsed++
	dbApp := &db.Application{
		Name: app.Title,
	}
	err := dbApp.Find()
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	dbApp.Usage++
	err = dbApp.Save()
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		ReloadAppList()
	}
}

func getAnalyticTools() *AppCategory {
	tools := &AppCategory{
		Category: "Analytic Tools",
	}

	kelly := newAppCart("kelly-criterion", "kelly_criterion.gohtml", "", "fas fa-coins",
		"/app/kelly-criterion", "Kelly Criterion", "Simulator for Kelly criterion. Kelly Criterion is a formula for sizing bets or investments from which the investor expects a positive return.")
	tools.AppCards = append(tools.AppCards, kelly)

	betSimulator := newAppCart("hilo-simulator", "hilo_simulator.gohtml", "", "fas fa-coins",
		"/app/hilo-simulator", "HiLo Simulator", "Simulate online hi/low betting website result(Provably fair only).")

	tools.AppCards = append(tools.AppCards, betSimulator)

	sortAppCardSlice(tools.AppCards)
	return tools
}

func getFormatTools() *AppCategory {
	tools := &AppCategory{
		Category: "Formatter",
	}

	encoderDecoder := newAppCart("string-encoder-decoder", "string_encoder_decoder.gohtml", "", "fas fa-receipt",
		"/app/string-encoder-decoder", "Encoder Decoder", "Convert string encoding based on request.")
	tools.AppCards = append(tools.AppCards, encoderDecoder)

	sortAppCardSlice(tools.AppCards)
	return tools
}

func getGeneratorTools() *AppCategory {
	tools := &AppCategory{
		Category: "Generator",
	}

	encoderDecoder := newAppCart("qr-code", "qr_code.gohtml", "", "fas fa-qrcode", "/app/qr-code", "QR Code", "Generate QR Code with animated logo and background image.")
	tools.AppCards = append(tools.AppCards, encoderDecoder)

	sortAppCardSlice(tools.AppCards)
	return tools
}

func getLookupTools() *AppCategory {
	tools := &AppCategory{
		Category: "Lookup",
	}

	dnsLookup := newAppCart("dns-lookup", "dns_lookup.gohtml", "", "fas fa-receipt",
		"/app/dns-lookup", "DNS Lookup", "Lookup given domain's DNS record (A, CNAME, PTR, NS, MX, TXT, and etc.")
	tools.AppCards = append(tools.AppCards, dnsLookup)

	ipLocation := newAppCart("ip-location", "ip_location.gohtml", "", "fas fa-map-marker-alt",
		"/app/ip-location", "IP Location", "Locate IP address to a specific location, by IPv4, or IPv6.")
	tools.AppCards = append(tools.AppCards, ipLocation)

	portChecker := newAppCart("port-checker", "port_checker.gohtml", "", "fas fa-network-wired", "/app/port-checker", "Port Checker", "Check if host port is open for tcp or udp pconnections.")
	tools.AppCards = append(tools.AppCards, portChecker)

	sortAppCardSlice(tools.AppCards)
	return tools
}

func getWebUtils() *AppCategory {
	tools := &AppCategory{
		Category: "Web Utils",
	}

	dnsLookup := newAppCart("request-bin", "request_bin.gohtml", "", "fas fa-receipt",
		"/app/request-bin", "Request Bin", "Receive and visualize HTTP requests for any method.")
	tools.AppCards = append(tools.AppCards, dnsLookup)

	sortAppCardSlice(tools.AppCards)
	return tools
}

func getPopular(apps []*AppCategory) *AppCategory {
	popular := &AppCategory{
		Category: "Popular",
	}

	for _, c := range apps {
		for _, a := range c.AppCards {
			popular.AppCards = append(popular.AppCards, a)
		}
	}

	sort.Slice(popular.AppCards, func(i, j int) bool {
		if popular.AppCards[i].AmountUsed != popular.AppCards[j].AmountUsed {
			return popular.AppCards[i].AmountUsed > popular.AppCards[j].AmountUsed
		}
		return popular.AppCards[i].Title > popular.AppCards[j].Title
	})

	return popular
}

func loadAppCardsUsage(appList []*AppCategory) {
	for _, category := range appList {
		for _, appCard := range category.AppCards {
			app := &db.Application{
				Name: appCard.Title,
			}
			err := app.Find()
			if err == gorm.ErrRecordNotFound {
				app.Name = appCard.Title
				app.Usage = 0
				app.Liked = 0
				err = app.Save()
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

func sortAppCardSlice(appCards []*AppCard) {
	sort.Slice(appCards, func(i, j int) bool {
		return strings.ToLower(appCards[i].Title) < strings.ToLower(appCards[j].Title)
	})
}

func newAppCart(name, templateName, imageURL, fontsAwesomeTag, link, title, description string) *AppCard {
	return &AppCard{
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
