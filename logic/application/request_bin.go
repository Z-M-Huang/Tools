package applicationlogic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	webData "github.com/Z-M-Huang/Tools/data/webdata/application"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/go-redis/redis"
)

//GetRequestBinHistory get request history by id
func GetRequestBinHistory(id string) *webData.RequestBinPageData {
	key := GetRequestBinKey(id)
	data := &webData.RequestBinPageData{}
	err := db.RedisGet(key, data)
	if err == redis.Nil {
		return nil
	} else if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return data
}

//NewRequestBinHistory new request history
func NewRequestBinHistory(private bool) *webData.RequestBinPageData {
	binData := &webData.RequestBinPageData{
		ID: strconv.FormatInt(time.Now().Unix(), 10),
	}
	if data.Config.HTTPS {
		binData.URL = "https://"
	} else {
		binData.URL = "http://"
	}

	binData.URL += data.Config.Host + "/api/request-bin/receive/" + binData.ID

	if private {
		binData.VerificationKey = utils.RandomString(30)
	}

	bytes, err := json.Marshal(binData)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	key := GetRequestBinKey(binData.ID)
	err = db.RedisSet(key, bytes, 24*time.Hour)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil
	}
	return binData
}

//GetRequestBinKey request bin key
func GetRequestBinKey(id string) string {
	return fmt.Sprintf("APP_REQUEST_BIN_%s", id)
}
