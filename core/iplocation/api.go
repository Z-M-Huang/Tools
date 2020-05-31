package iplocation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//API iplocation
type API struct{}

var client *http.Client
var rateLimitKey string = "APP_IP_LOCATION_RATE_LIMIT"

func init() {
	client = &http.Client{}
	client.Timeout = 1 * time.Second
}

//Get /api/ip-location/get
func (API) Get(c *gin.Context) {
	if !db.RedisExist(rateLimitKey) {
		err := db.RedisSet(rateLimitKey, 1, 24*time.Hour)
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
	} else {
		num, err := db.RedisGetInt(rateLimitKey)
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		if num >= 950 {
			c.Status(http.StatusTooManyRequests)
			return
		}
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://ipapi.co/%s/json/", c.ClientIP()), nil)
	req.Header.Add("DNT", "1")
	req.Header.Add("User-Agent", "https://tools.zh-code.com")
	resp, err := client.Do(req)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	db.RedisIncr(rateLimitKey)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		utils.Logger.Sugar().Errorf("https://ipapi.co/%s/json/ returned %s", c.ClientIP(), resp.Status)
		c.Status(http.StatusServiceUnavailable)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	response := &Response{}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		utils.Logger.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response)
}
