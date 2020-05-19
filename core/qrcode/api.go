package qrcode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"strconv"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/Z-M-Huang/go-qrcode"
	"github.com/gin-gonic/gin"
)

//API qrcode
type API struct{}

//CreateQRCode /api/qr-code/create
func (API) CreateQRCode(c *gin.Context) {
	response := core.GetResponseInContext(c.Keys)
	c.Request.ParseMultipartForm(1024)
	request := &Request{
		Content:         c.Request.PostFormValue("content"),
		Level:           c.Request.PostFormValue("level"),
		BackgroundColor: c.Request.PostFormValue("backColor"),
		ForegroundColor: c.Request.PostFormValue("foreColor"),
	}
	var err error
	request.Size, err = strconv.Atoi(c.Request.PostFormValue("size"))
	if err != nil {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Invalid request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	if len(request.Content) < 1 {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Invalid request.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	if request.Size > 1024 {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Invalid request. The size is too big.",
		})
		core.WriteResponse(c, 400, response)
		return
	} else if request.Size < 0 {
		response.SetAlert(&core.AlertData{
			IsWarning: true,
			Message:   "Invalid request. The size cannot be negative",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	var level qrcode.RecoveryLevel
	switch request.Level {
	case "L":
		level = qrcode.Low
	case "M":
		level = qrcode.Medium
	case "Q":
		level = qrcode.High
	case "H":
		level = qrcode.Highest
	default:
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Invalid request. Invalid Level.",
		})
		core.WriteResponse(c, 400, response)
		return
	}

	var backgroundColor color.Color = color.White
	if request.BackgroundColor != "" {
		backgroundColor, err = parseHexColorFast(request.BackgroundColor)
		if err != nil {
			response.SetAlert(&core.AlertData{
				IsWarning: true,
				Message:   "Background Color: " + err.Error(),
			})
			core.WriteResponse(c, 400, response)
			return
		}
	}

	var foregroundColor color.Color = color.Black
	if request.ForegroundColor != "" {
		foregroundColor, err = parseHexColorFast(request.ForegroundColor)
		if err != nil {
			response.SetAlert(&core.AlertData{
				IsWarning: true,
				Message:   "Foreground Color: " + err.Error(),
			})
			core.WriteResponse(c, 400, response)
			return
		}
	}

	var logo image.Image
	logoFile, _, err := c.Request.FormFile("logoImage")
	if err == nil {
		logo, _, err = image.Decode(logoFile)
		if err != nil {
			response.SetAlert(&core.AlertData{
				IsWarning: true,
				Message:   "Failed to get logo image",
			})
			core.WriteResponse(c, 400, response)
			return
		}
	}

	var backgroundImage image.Image
	backgroundImageFile, _, err := c.Request.FormFile("backgroundImage")
	if err == nil {
		backgroundImage, _, err = image.Decode(backgroundImageFile)
		if err != nil {
			response.SetAlert(&core.AlertData{
				IsWarning: true,
				Message:   "Failed to get background image",
			})
			core.WriteResponse(c, 400, response)
			return
		}
	}

	q, err := qrcode.New(request.Content, level)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Internal Error",
		})
		core.WriteResponse(c, 500, response)
		return
	}

	q.BackgroundColor = backgroundColor
	q.ForegroundColor = foregroundColor
	q.LogoImage = logo
	q.BackgroundImage = backgroundImage

	redisKey := getRedisKey(c.ClientIP())
	if !db.RedisExist(redisKey) {
		db.RedisSet(redisKey, 19, 24*time.Hour)
	} else {
		val, err := db.RedisDecr(redisKey)
		if err != nil {
			response.SetAlert(&core.AlertData{
				IsDanger: true,
				Message:  "Internal Error",
			})
			core.WriteResponse(c, 500, response)
			return
		} else if val < 0 {
			response.SetAlert(&core.AlertData{
				IsWarning: true,
				Message:   "Too many requests today. Please come back tomorrow.",
			})
			core.WriteResponse(c, http.StatusTooManyRequests, response)
			return
		}
	}

	imageBytes, err := q.PNG(request.Size)
	if err != nil {
		utils.Logger.Error(err.Error())
		response.SetAlert(&core.AlertData{
			IsDanger: true,
			Message:  "Internal Error",
		})
		core.WriteResponse(c, 500, response)
		return
	}

	response.Data = base64.StdEncoding.EncodeToString(imageBytes)
	core.WriteResponse(c, 200, response)
}

func getRedisKey(ip string) string {
	return fmt.Sprintf("APP_QR_CODE_%s", ip)
}

func parseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		s = "#" + s
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errors.New("Invalid color format")
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errors.New("Invalid color format")
	}
	return
}
