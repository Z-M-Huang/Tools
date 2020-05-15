package applicationlogic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	webData "github.com/Z-M-Huang/Tools/data/webdata/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

//GetRequestBinHistory get request history by id
func GetRequestBinHistory(c *gin.Context, id string) *webData.RequestBinPageData {
	key := getRequestKey(id)
	val, err := utils.RedisClient.Get(key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	data := &webData.RequestBinPageData{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return data
}

//NewRequestBinHistory new request history
func NewRequestBinHistory(c *gin.Context, private bool) *webData.RequestBinPageData {
	data := &webData.RequestBinPageData{
		ID: strconv.FormatInt(time.Now().Unix(), 10),
	}
	if utils.Config.HTTPS {
		data.URL = "https://"
	} else {
		data.URL = "http://"
	}

	data.URL += utils.Config.Host + "/app/request-bin/" + data.ID

	if private {
		data.VerificationKey = utils.RandomString(30)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	key := getRequestKey(data.ID)
	err = utils.RedisClient.Set(key, bytes, 7*24*time.Hour).Err()
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return data
}

func getRequestKey(id string) string {
	return fmt.Sprintf("APP_REQUEST_BIN_%s", id)
}
