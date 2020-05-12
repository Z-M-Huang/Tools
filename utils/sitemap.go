package utils

import (
	"fmt"

	"github.com/ikeikeikeike/go-sitemap-generator/stm"
)

//BuildSitemap build sitemap
func BuildSitemap() *stm.Sitemap {
	sm := stm.NewSitemap()
	host := ""
	if Config.SitemapConfig.IsHTTPS {
		host = "https://"
	} else {
		host = "http://"
	}
	sm.SetDefaultHost(fmt.Sprintf("%s%s", host, Config.Host))
	sm.SetCompress(true)

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
