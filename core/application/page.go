package application

import (
	"net/http"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/core/requestbin"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/blevesearch/bleve"
	"github.com/gin-gonic/gin"
)

//Page logic
type Page struct{}

//RenderApplicationPage renders /app/:name
func (Page) RenderApplicationPage(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.PageResponse)

	name := c.Param("name")

	if name == "" {
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}

	appCard := GetApplicationsByName(name)
	if appCard == nil {
		c.String(http.StatusNotFound, "404 Not Found")
		return
	}
	response.Header.Title = appCard.Title + " - Fun Apps"
	response.Header.Description = appCard.Description

	AddApplicationUsage(appCard)
	response.Data = loadAppSpecificData(c, appCard.Name)
	if c.IsAborted() {
		return
	}

	c.HTML(http.StatusOK, appCard.TemplateName, response)
}

//SearchApps search apps in search bar
func (Page) SearchApps(c *gin.Context) {
	response := c.Keys[utils.ResponseCtxKey].(*data.PageResponse)
	claim := account.GetClaimInContext(c.Keys)
	keys := c.Query("keywords")
	if keys == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	if searchIndex == nil {
		core.WriteUnexpectedError(c, response)
		return
	}
	query := bleve.NewQueryStringQuery(keys)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := searchIndex.Search(searchRequest)
	if err != nil {
		core.WriteUnexpectedError(c, response)
		return
	}
	if len(searchResult.Hits) < 1 {
		c.HTML(http.StatusOK, "app_search.gohtml", response)
		return
	}
	var names []string
	for _, h := range searchResult.Hits {
		names = append(names, h.ID)
	}
	if !(claim == nil) {
		user := &db.User{
			Email: claim.Id,
		}
		err := user.Find()
		if err == nil && len(user.LikedApps) > 0 {
			response.Data = SearchAppListByNamesWithLikes(user, names)
		}
	} else {
		response.Data = SearchAppListByNames(names)
	}
	response.Header.Title = "Search - Fun Apps"
	response.Header.Description = "Fun apps, fun personal small projects, and just for fun."
	c.HTML(http.StatusOK, "app_search.gohtml", response)
}

func loadAppSpecificData(c *gin.Context, appName string) interface{} {
	switch appName {
	case "request-bin":
		return requestbin.LoadRequestBinData(c)
	case "port-checker":
		return c.ClientIP()
	}
	return nil
}
