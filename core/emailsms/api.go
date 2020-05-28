package emailsms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/gin-gonic/gin"
)

//API emailsms
type API struct{}

func sendEmail(toAddress, subject, content string) error {
	from := mail.Address{
		Name:    "",
		Address: data.Config.EmailConfig.EmailAddress,
	}
	to := mail.Address{
		Name:    "",
		Address: toAddress,
	}
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content

	//smtp.gmail.com:465
	host, _, _ := net.SplitHostPort(data.Config.EmailConfig.SMTPServer)

	auth := smtp.PlainAuth("", data.Config.EmailConfig.EmailAddress, data.Config.EmailConfig.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", data.Config.EmailConfig.SMTPServer, tlsconfig)
	if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("Failed to connect to email server")
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("InternalServer Error")
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("InternalServer Error")
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("InternalServer Error")
	}

	if err = c.Rcpt(to.Address); err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("Failed to set to address")
	}

	// Data
	w, err := c.Data()
	if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("Failed to write message content")
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("Failed to write message content")
	}

	err = w.Close()
	if err != nil {
		utils.Logger.Error(err.Error())
		return errors.New("InternalServer Error")
	}

	c.Quit()
	return nil
}

//Send /api/emailsms/send
func (API) Send(c *gin.Context) {
	response := &data.APIResponse{}
	request := &Request{}
	err := c.ShouldBind(&request)
	if err != nil {
		response.Message = "Bad Request"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}
}

func getCarrierGateway(carrier string) string {
	switch carrier {
	case "AT&T":
		return "@mms.att.net"
	case "T-Mobile":
		return "@tmomail.net"
	case "Verizon":
		return "@vzwpix.com"
	case "Sprint":
		return "@pm.sprint.com"
	case "Xfinity":
		return "@mypixmessages.com"
	case "Virgin Mobile":
		return "@vmpix.com"
	case "Tracfone":
		return "@mmst5.tracfone.com"
	case "Simple Mobile":
		return "@smtext.com"
	case "Mint Mobile":
		return "@mailmymobile.net"
	case "Red Pocket", "Page Plus":
		return "@vtext.com"
	case "Metro PCS":
		return "@mymetropcs.com"
	case "Boost Mobile":
		return "@myboostmobile.com"
	case "Cricket":
		return "@mms.cricketwireless.net"
	case "Republic Wireless":
		return "@text.republicwireless.com"
	case "Google Fi":
		return "@msg.fi.google.com"
	case "U.S. Cellular":
		return "@mms.uscc.net"
	case "Ting":
		return "@message.ting.com"
	case "Consumer Cellular":
		return "@mailmymobile.net"
	case "C-Spire":
		return "@cspire1.com"
	}
	return ""
}
