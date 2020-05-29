package qrcode

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"net/http"
	"strconv"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/Z-M-Huang/go-qrcode"
	"github.com/gin-gonic/gin"
)

//API qrcode
type API struct{}

// CreateQRCode /api/qr-code/create
// @Summary Generate QR Code with animated logo and background image.
// @Description Generate QR Code with animated logo and background image.
// @Tags Generator
// @Accept mpfd
// @Produce json,xml
// @Param content body string true "Content in QRCode"
// @Param level body string true "Level L|M|Q|H"
// @Param size body int true "Size of the image"
// @Param backColor body string false "Hex color"
// @Param foreColor body string false "Hex color"
// @Param logoImage formData file false "Logo image file"
// @Param logoGifImage formData file false "Logo gif image file"
// @Param backgroundImage formData file false "Background image file"
// @Success 200 {object} data.APIResponse
// @Failure 400 {object} data.APIResponse
// @Failure 429 {object} data.APIResponse
// @Failure 500 {object} data.APIResponse
// @Router /api/qr-code/create [post]
func (API) CreateQRCode(c *gin.Context) {
	response := &data.APIResponse{}
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
		response.Message = "Invalid Request."
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	if len(request.Content) < 1 {
		response.Message = "Invalid Request."
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	if request.Size > 1024 {
		response.Message = "Invalid Request. The size is too big."
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	} else if request.Size < 0 {
		response.Message = "Invalid Request. The size cannot be negative."
		core.WriteResponse(c, http.StatusBadRequest, response)
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
		response.Message = "Invalid Request. Invalid Level."
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	var backgroundColor color.Color = color.White
	if request.BackgroundColor != "" {
		backgroundColor, err = parseHexColorFast(request.BackgroundColor)
		if err != nil {
			response.Message = "Background Color: " + err.Error()
			core.WriteResponse(c, http.StatusBadRequest, response)
			return
		}
	}

	var foregroundColor color.Color = color.Black
	if request.ForegroundColor != "" {
		foregroundColor, err = parseHexColorFast(request.ForegroundColor)
		if err != nil {
			response.Message = "Foreground Color: " + err.Error()
			core.WriteResponse(c, http.StatusBadRequest, response)
			return
		}
	}

	var logo image.Image
	logoFile, _, err := c.Request.FormFile("logoImage")
	if err == nil {
		logo, _, err = image.Decode(logoFile)
		if err != nil {
			response.Message = "Failed to get logo image."
			core.WriteResponse(c, http.StatusBadRequest, response)
			return
		}
	}

	var logoGifImage *gif.GIF
	logoGifFile, _, err := c.Request.FormFile("logoGifImage")
	if err == nil {
		logoGifImage, err = gif.DecodeAll(logoGifFile)
		if err != nil {
			response.Message = "Failed to get logo gif image."
			core.WriteResponse(c, http.StatusBadRequest, response)
			return
		}
	}

	var backgroundImage image.Image
	backgroundImageFile, _, err := c.Request.FormFile("backgroundImage")
	if err == nil {
		backgroundImage, _, err = image.Decode(backgroundImageFile)
		if err != nil {
			response.Message = "Failed to get background image."
			core.WriteResponse(c, http.StatusBadRequest, response)
			return
		}
	}

	q, err := qrcode.New(request.Content, level)
	if err != nil {
		response.Message = "Internal Error"
		core.WriteResponse(c, http.StatusInternalServerError, response)
		return
	}

	q.BackgroundColor = backgroundColor
	q.ForegroundColor = foregroundColor
	q.BackgroundImage = backgroundImage

	redisKey := getRedisKey(c.ClientIP())
	if !db.RedisExist(redisKey) {
		db.RedisSet(redisKey, 19, 24*time.Hour)
	} else {
		val, err := db.RedisDecr(redisKey)
		if err != nil {
			response.Message = "Internal Error"
			core.WriteResponse(c, http.StatusInternalServerError, response)
			return
		} else if val < 0 {
			response.Message = "Too many requests today. Please come back tomorrow."
			core.WriteResponse(c, http.StatusTooManyRequests, response)
			return
		}
	}

	var imageBytes []byte
	if logo != nil {
		imageBytes, err = q.PNGWithLogo(request.Size, logo)
	} else if logoGifImage != nil {
		imageBytes, err = q.GIFLogo(request.Size, logoGifImage)
	} else {
		imageBytes, err = q.PNG(request.Size)
	}
	if err != nil {
		utils.Logger.Error(err.Error())
		response.Message = "Internal Error"
		core.WriteResponse(c, http.StatusInternalServerError, response)
		return
	}

	response.Data = base64.StdEncoding.EncodeToString(imageBytes)
	core.WriteResponse(c, http.StatusOK, response)
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
