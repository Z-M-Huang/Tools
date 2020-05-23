package qrcode

import "mime/multipart"

//Request /api/qr-code/create
type Request struct {
	Content         string                `json:"content" xml:"content" form:"content" bind:"required"`
	Level           string                `json:"level" xml:"level" form:"level" bind:"required"`
	Size            int                   `json:"size" xml:"size" form:"size"`
	BackgroundColor string                `json:"backColor" xml:"backColor" form:"backColor"`
	ForegroundColor string                `json:"foreColor" xml:"foreColor" form:"foreColor"`
	LogoImage       *multipart.FileHeader `json:"logoImage" xml:"logoImage" form:"logoImage"`
	LogoGifImage    *multipart.FileHeader `json:"logoGifImage" xml:"logoGifImage" form:"logoGifImage"`
	BackgroundImage *multipart.FileHeader `json:"backgroundImage" xml:"backgroundImage" form:"backgroundImage"`
}
