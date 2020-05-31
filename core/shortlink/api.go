package shortlink

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/core/account"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var lock sync.Mutex

//API shortlink API
type API struct{}

// Get /api/shortlink/get
// @Summary Get short link
// @Description Get short link. Rate limit: 1,000 per hour
// @Tags Short-link
// @Accept json
// @Produce json,xml
// @Param "" body Request true "Request JSON"
// @Success 200 {object} data.APIResponse
// @Failure 400 {object} data.APIResponse
// @Failure 429 {object} data.APIResponse
// @Failure 500 {object} data.APIResponse
// @Router /api/shortlink/get [post]
func (API) Get(c *gin.Context) {
	response := &data.APIResponse{}
	request := &Request{}

	err := c.ShouldBind(&request)
	if err != nil {
		response.Message = "Bad Request"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	request.URL = strings.TrimSpace(request.URL)
	_, err = url.Parse(request.URL)
	if err != nil {
		response.Message = "Bad Request"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	ipKey := getIPKey(c.ClientIP())
	if !db.RedisExist(ipKey) {
		db.RedisSet(ipKey, 1, 1*time.Hour)
	} else {
		counter, err := db.RedisGetInt(ipKey)
		if err != nil {
			response.Message = "InternalServer Error. Please try again..."
			core.WriteResponse(c, http.StatusInternalServerError, response)
			return
		}
		if counter > 1000 {
			response.Message = "Too many requests. Rate limit is 1000 per hour."
			core.WriteResponse(c, http.StatusTooManyRequests, response)
			return
		}
		db.RedisIncr(ipKey)
	}

	shortLink := &db.ShortLink{
		Link: &request.URL,
	}

	claim := account.GetClaimInContext(c.Keys)
	if claim != nil {
		user := &db.User{
			Email: claim.Id,
		}
		err = user.Find()
		if err != nil {
			utils.Logger.Error(err.Error())
		} else {
			shortLink.UserID = user.ID
		}
	}

	err = shortLink.Save()
	if err != nil {
		response.Message = "Failed to get a new link. Please try again..."
		core.WriteResponse(c, http.StatusInternalServerError, response)
		return
	}

	host := ""
	if data.Config.HTTPS {
		host = "https://" + data.Config.Host
	} else {
		host = "http://" + data.Config.Host
	}
	response.Data = fmt.Sprintf("%s/s/%d", host, shortLink.ID)
	core.WriteResponse(c, http.StatusOK, response)
}

// RedirectShortLink /s/:id
// @Summary Redrect to short link
// @Description Redrect to short link
// @Tags Short-link
// @Accept json
// @Produce json,xml
// @Param id path int true "Short Link ID"
// @Success 200 {object} data.APIResponse
// @Failure 400 {object} data.APIResponse
// @Failure 500 {object} data.APIResponse
// @Router /s/{id} [post]
func (API) RedirectShortLink(c *gin.Context) {
	response := &data.APIResponse{}
	idStr := c.Param("id")
	if idStr == "" {
		response.Message = "Not found"
		core.WriteResponse(c, http.StatusNotFound, response)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Message = "Bad Request"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	shortlink := &db.ShortLink{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	lock.Lock()
	defer lock.Unlock()
	if shortlink.Find() != nil {
		response.Message = "Not Found"
		core.WriteResponse(c, http.StatusNotFound, response)
		return
	}

	shortlink.Usage++
	shortlink.Save()

	c.Redirect(http.StatusPermanentRedirect, *shortlink.Link)
}

func getIPKey(ip string) string {
	return fmt.Sprintf("APP_SHORT_LINK_IP_%s", ip)
}
