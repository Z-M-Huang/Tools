package application

import "mime/multipart"

//QRCodeRequest /api/qr-code/create
type QRCodeRequest struct {
	Content         string                `json:"content" xml:"content" form:"content" bind:"required"`
	Level           string                `json:"level" xml:"level" form:"level"`
	Size            int                   `json:"size" xml:"size" form:"size"`
	BackgroundColor string                `json:"backColor" xml:"backColor" form:"backColor"`
	ForegroundColor string                `json:"foreColor" xml:"foreColor" form:"foreColor"`
	LogoImage       *multipart.FileHeader `json:"logoImage" xml:"logoImage" form:"logoImage"`
	BackgroundImage *multipart.FileHeader `json:"backgroundImage" xml:"backgroundImage" form:"backgroundImage"`
}
