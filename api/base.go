package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Z-M-Huang/Tools/data"
	"github.com/gin-gonic/gin"
)

//WriteResponse Write api response
func WriteResponse(c *gin.Context, code int, obj interface{}) {
	switch c.NegotiateFormat(gin.MIMEJSON, gin.MIMEXML, gin.MIMEYAML, gin.MIMEPlain) {
	case gin.MIMEJSON:
		c.JSON(code, obj)
	case gin.MIMEXML:
		c.XML(code, obj)
	case gin.MIMEYAML:
		c.YAML(code, obj)
	default:
		c.String(http.StatusNotAcceptable, "Not Acceptable")
	}
}

//WriteUnexpectedError Write unexpected api response
func WriteUnexpectedError(c *gin.Context, response *data.Response) {
	response.SetAlert(&data.AlertData{
		IsDanger: true,
		Message:  "Um... Your data got eaten by the cyber space... Would you like to try again?",
	})
	WriteResponse(c, http.StatusInternalServerError, response)
}

//ParseFloat parse float from string
func ParseFloat(input string, size int, required bool) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		if required {
			return 0, errors.New("Cannot be empty")
		}
		return 0, nil
	}

	ret, err := strconv.ParseFloat(input, size)
	if err != nil {
		if required {
			return 0, errors.New("Invalid input, needs to be a number")
		}
		return 0, nil
	}
	return ret, nil
}

//ParseUint parse uint from string
func ParseUint(input string, base, size int, required bool) (uint64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		if required {
			return 0, errors.New("Cannot be empty")
		}
		return 0, nil
	}

	ret, err := strconv.ParseUint(input, base, size)
	if err != nil {
		if required {
			return 0, errors.New("Invalid input, needs to be an integer")
		}
		return 0, nil
	}
	return ret, nil
}

//ParseBool parse boolean from string
func ParseBool(input string, required bool) (bool, error) {
	input = strings.TrimSpace(input)
	if input == "on" {
		return true, nil
	}

	ret, err := strconv.ParseBool(input)
	if err != nil {
		if required {
			return false, errors.New("Invalid input, needs to be a boolean")
		}
		return false, nil
	}
	return ret, nil
}
