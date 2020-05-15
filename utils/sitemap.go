package utils

import (
	"fmt"
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

var client http.Client

//BuildSitemap build sitemap
func BuildSitemap() *stm.Sitemap {
	sm := stm.NewSitemap()
	host := ""
	if Config.IsHTTPS {
		host = "https://"
	} else {
		host = "http://"
	}
	sm.SetDefaultHost(fmt.Sprintf("%s%s", host, Config.Host))

	sm.Create()

	sm.Add(getPageSiteMap("/"))
	sm.Add(getPageSiteMap("/login"))
	sm.Add(getPageSiteMap("/signup"))

	for _, category := range AppList {
		for _, app := range category.AppCards {
			sm.Add(getPageSiteMap(app.Link))
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
		"Disallow: /api/*"}

	sitemapURL := ""
	if Config.IsHTTPS {
		sitemapURL = fmt.Sprintf("https://%s/sitemap.xml", Config.Host)
	} else {
		sitemapURL = fmt.Sprintf("http://%s/sitemap.xml", Config.Host)
	}
	content = append(content, fmt.Sprintf("Sitemap: %s", sitemapURL))
	return content
}

//PingSearchEngines ping search engine
func PingSearchEngines() {
	sitemapURL := ""
	if Config.IsHTTPS {
		sitemapURL = fmt.Sprintf("https://%s/sitemap.xml", Config.Host)
	} else {
		sitemapURL = fmt.Sprintf("http://%s/sitemap.xml", Config.Host)
	}

	resp, err := http.Get(fmt.Sprintf("http://www.google.com/webmasters/tools/ping?sitemap=%s", sitemapURL))
	if err != nil {
		Logger.Sugar().Errorf("Failed to ping Google Search Engine %s", err.Error())
	} else if resp.StatusCode != 200 {
		Logger.Sugar().Errorf("Failed to ping Google Search Engine. Returned %s", resp.Status)
		resp.Body.Close()
	} else {
		Logger.Info("Successfully Pinged Google")
		resp.Body.Close()
	}

	resp, err = http.Get(fmt.Sprintf("http://www.bing.com/webmaster/ping.aspx?siteMap=%s", sitemapURL))
	if err != nil {
		Logger.Sugar().Errorf("Failed to ping Bing Search Engine %s", err.Error())
	} else if resp.StatusCode != 200 {
		Logger.Sugar().Errorf("Failed to ping Bing Search Engine. Returned %s", resp.Status)
		resp.Body.Close()
	} else {
		Logger.Info("Successfully Pinged Bing")
		resp.Body.Close()
	}
}
