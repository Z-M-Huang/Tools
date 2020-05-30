package main

import (
	"fmt"
	"net/http"

	"github.com/Z-M-Huang/Tools/core/application"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

var client http.Client

//BuildSitemap build sitemap
func BuildSitemap() *stm.Sitemap {
	sm := stm.NewSitemap()
	host := ""
	if data.Config.HTTPS {
		host = "https://"
	} else {
		host = "http://"
	}
	sm.SetDefaultHost(fmt.Sprintf("%s%s", host, data.Config.Host))

	sm.Create()

	sm.Add(getPageSiteMap("/"))
	sm.Add(getPageSiteMap("/login"))
	sm.Add(getPageSiteMap("/signup"))
	sm.Add(getPageSiteMap("/swagger/index.html"))

	for _, category := range application.GetAppList() {
		if category.Category != "Popular" {
			for _, app := range category.AppCards {
				sm.Add(getPageSiteMap(app.Link))
			}
		}
	}

	// Note: Do not call `sm.Finalize()` because it flushes
	// the underlying data structure from memory to disk.

	return sm
}

func getPageSiteMap(loc string) stm.URL {
	dic := stm.URL{}
	dic["loc"] = loc
	dic["changefreq"] = "daily"
	dic["mobile"] = "true"
	return dic
}

//GetRobotsTxt get robots.txt content
func GetRobotsTxt() []string {
	content := []string{"User-agent: *",
		"Disallow: /api/*",
		"Disallow: /login",
		"Disallow: /signup",
		"Disallow: /swagger/index.html"}

	sitemapURL := ""
	if data.Config.HTTPS {
		sitemapURL = fmt.Sprintf("https://%s/sitemap.xml", data.Config.Host)
	} else {
		sitemapURL = fmt.Sprintf("http://%s/sitemap.xml", data.Config.Host)
	}
	content = append(content, fmt.Sprintf("Sitemap: %s", sitemapURL))
	return content
}

//PingSearchEngines ping search engine
func PingSearchEngines() {
	sitemapURL := ""
	if data.Config.HTTPS {
		sitemapURL = fmt.Sprintf("https://%s/sitemap.xml", data.Config.Host)
	} else {
		sitemapURL = fmt.Sprintf("http://%s/sitemap.xml", data.Config.Host)
	}

	resp, err := http.Get(fmt.Sprintf("http://www.google.com/webmasters/tools/ping?sitemap=%s", sitemapURL))
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to ping Google Search Engine %s", err.Error())
	} else if resp.StatusCode != 200 {
		utils.Logger.Sugar().Errorf("Failed to ping Google Search Engine. Returned %s", resp.Status)
		resp.Body.Close()
	} else {
		utils.Logger.Info("Successfully Pinged Google")
		resp.Body.Close()
	}

	resp, err = http.Get(fmt.Sprintf("http://www.bing.com/webmaster/ping.aspx?siteMap=%s", sitemapURL))
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to ping Bing Search Engine %s", err.Error())
	} else if resp.StatusCode != 200 {
		utils.Logger.Sugar().Errorf("Failed to ping Bing Search Engine. Returned %s", resp.Status)
		resp.Body.Close()
	} else {
		utils.Logger.Info("Successfully Pinged Bing")
		resp.Body.Close()
	}
}
