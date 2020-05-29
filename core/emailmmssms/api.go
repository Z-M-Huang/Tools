package emailmmssms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/Z-M-Huang/Tools/core"
	"github.com/Z-M-Huang/Tools/data"
	"github.com/Z-M-Huang/Tools/data/db"
	"github.com/Z-M-Huang/Tools/utils"
	"github.com/emersion/go-imap"
	idle "github.com/emersion/go-imap-idle"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/gin-gonic/gin"
)

//API emailsms
type API struct{}

var phoneRe = regexp.MustCompile(`^\d{10}$`)
var maxDailyEmailAmount int64 = 1500

const redisLockKey string = "LOCK_EMAIL_MMS_SMS"

func sendEmail(toAddress, subject, content, ipAddress string) error {
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
	headers["Client-IP"] = ipAddress

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

//Send /api/email-mms-sms/send
func (API) Send(c *gin.Context) {
	response := &data.APIResponse{}
	request := &Request{}
	err := c.ShouldBind(&request)
	if err != nil {
		response.Message = "Bad Request"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	if !db.RedisExist(getTotalEmailKey()) {
		err = db.RedisSet(getTotalEmailKey(), 1, 24*time.Hour)
		if err != nil {
			utils.Logger.Error(err.Error())
			response.Message = "InternalServer Error. Please try again later"
			core.WriteResponse(c, http.StatusInternalServerError, response)
			return
		}
	} else {
		currentCount, err := db.RedisGetInt(getTotalEmailKey())
		if err != nil {
			utils.Logger.Error(err.Error())
			response.Message = "InternalServer Error. Please try again later"
			core.WriteResponse(c, http.StatusInternalServerError, response)
			return
		}
		if currentCount >= maxDailyEmailAmount {
			utils.Logger.Info("Server exceed usage limit. Please submit a ticket to report this issue.")
			response.Message = "Server exceed usage limit. Please submit a ticket to report this issue."
			core.WriteResponse(c, http.StatusInternalServerError, response)
			return
		}
		db.RedisIncr(getTotalEmailKey())
	}

	request.Subject = strings.TrimSpace(request.Subject)
	request.Content = strings.TrimSpace(request.Content)
	request.ToNumber = strings.TrimSpace(request.ToNumber)
	if !phoneRe.Match([]byte(request.ToNumber)) {
		response.Message = "Bad Request. Only 10 digit US or Canada phone number supported."
		core.WriteResponse(c, http.StatusBadRequest, response)
	}

	if request.Subject == "" && request.Content == "" {
		response.Message = "Bad Request. Content unknown."
		core.WriteResponse(c, http.StatusBadRequest, response)
	}

	carrier := getCarrierGateway(strings.TrimSpace(request.Carrier))
	if carrier == "" {
		response.Message = "Bad Request.Unknown Carrier"
		core.WriteResponse(c, http.StatusBadRequest, response)
		return
	}

	status, err := checkIPToNumberAllowed(c.ClientIP(), request.ToNumber)
	if status != http.StatusOK || err != nil {
		response.Message = err.Error()
		core.WriteResponse(c, status, response)
		return
	}

	err = sendEmail(fmt.Sprintf("%s%s", request.ToNumber, carrier), request.Subject, request.Content, c.ClientIP())
	if err != nil {
		response.Message = err.Error()
		core.WriteResponse(c, http.StatusInternalServerError, response)
		return
	}
	response.Data = true
	core.WriteResponse(c, http.StatusOK, response)
	return
}

//StartListeningToEmails imap client listening for new emails
func (API) StartListeningToEmails() {
	if data.Config.EmailConfig.IMAPServer == "" {
		return
	}

	lock := db.RedisLock(redisLockKey, 1*time.Second)
	if lock != nil {
		c, err := imapClient.DialTLS(data.Config.EmailConfig.IMAPServer, nil)
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}
		utils.Logger.Info("IMAP connected")

		if err := c.Login(data.Config.EmailConfig.EmailAddress, data.Config.EmailConfig.Password); err != nil {
			utils.Logger.Error(err.Error())
			return
		}
		utils.Logger.Info("IMAP Logged in")

		//Search for unseen emails first
		_, err = c.Select("INBOX", false)
		if err != nil {
			utils.Logger.Fatal(err.Error())
			return
		}

		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}
		ids, err := c.Search(criteria)
		if err != nil {
			utils.Logger.Fatal(err.Error())
			return
		}

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)

		//Load not seen messages
		if len(ids) > 0 {
			seqset := new(imap.SeqSet)
			seqset.AddNum(ids...)
			section := &imap.BodySectionName{}
			go func() {
				done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, section.FetchItem()}, messages)
			}()

			for msg := range messages {
				loadReplyToRedis(msg)
			}
		}

		idleClient := idle.NewClient(c)

		updates := make(chan imapClient.Update)
		c.Updates = updates
		// Start idling
		go func() {
			done <- idleClient.IdleWithFallback(nil, 0)
		}()

		// Listen for updates
		for {
			select {
			case update := <-updates:
				switch update.(type) {
				case *imapClient.MailboxUpdate:
					msg := update.(*imapClient.MailboxUpdate)
					uid := msg.Mailbox.Messages
					seqset := new(imap.SeqSet)
					seqset.AddRange(uid, uid)
					section := &imap.BodySectionName{}
					go func() {
						done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody, section.FetchItem()}, messages)
					}()
					for msg := range messages {
						loadReplyToRedis(msg)
					}
				default:
				}
			case err := <-done:
				if err != nil {
					utils.Logger.Fatal(err.Error())
				}
				utils.Logger.Fatal("Not idling anymore")
				return
			}
		}
	}
}

func loadReplyToRedis(msg *imap.Message) {
	r := msg.GetBody(&imap.BodySectionName{})
	mr, _ := mail.CreateReader(r)
	for {
		//Only search for address contains +
		if strings.Contains(msg.Envelope.To[0].Address(), "+") {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				utils.Logger.Error(err.Error())
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				if strings.Contains(p.Header.Get("Content-Type"), "text/plain") {
					b, _ := ioutil.ReadAll(p.Body)
					var histories []*MessageHistory
					key := getReplyMessageKeyByEmail(msg.Envelope.To[0].Address())
					if !db.RedisExist(key) {
						histories = append(histories, &MessageHistory{
							DateReceived: msg.Envelope.Date,
							Content:      string(b),
						})
						db.RedisSet(key, histories, 24*time.Hour)
					} else {
						err = db.RedisGet(key, histories)
						if err != nil {
							utils.Logger.Error(err.Error())
							continue
						}
						histories = append(histories, &MessageHistory{
							DateReceived: msg.Envelope.Date,
							Content:      string(b),
						})
						db.RedisSet(key, histories, 24*time.Hour)
					}
					utils.Logger.Sugar().Infof("Received Email MMS/SMS reply from %s", msg.Envelope.To[0].Address)
				}
			case *mail.AttachmentHeader:
				// This is an attachment do nothing
			}
		}
	}
}

func checkIPToNumberAllowed(ipAddress, toNumber string) (int, error) {
	ipKey := getRedisIPKey(ipAddress)
	toNumberKey := getRedisToNumberKey(toNumber)
	if !db.RedisExist(ipKey) {
		err := db.RedisSet(ipKey, 1, 1*time.Minute)
		if err != nil {
			utils.Logger.Error(err.Error())
			return http.StatusInternalServerError, errors.New("InternalServer Error. Please try again later")
		}
	} else {
		ipAmount, err := db.RedisGetInt(ipKey)
		if err != nil {
			return http.StatusInternalServerError, errors.New("InternalServer Error. Please try again later")
		}
		// 5 requests per IP
		if ipAmount > 4 {
			return http.StatusTooManyRequests, errors.New("Too many requests, 5 messages per IP per minute only")
		}
		db.RedisIncr(ipKey)
	}

	if !db.RedisExist(toNumberKey) {
		err := db.RedisSet(toNumberKey, 1, 1*time.Minute)
		if err != nil {
			utils.Logger.Error(err.Error())
			return http.StatusInternalServerError, errors.New("InternalServer Error. Please try again later")
		}
	} else {
		toNumberAmount, err := db.RedisGetInt(toNumberKey)
		if err != nil {
			return http.StatusInternalServerError, errors.New("InternalServer Error. Please try again later")
		}
		if toNumberAmount > 1 {
			return http.StatusTooManyRequests, errors.New("Too many requests to the specific number, 2 messages to specific phone number per minute only")
		}
		db.RedisIncr(toNumberKey)
	}

	return http.StatusOK, nil
}

func getRedisIPKey(ip string) string {
	return fmt.Sprintf("APP_EMAIL_SMS_IP_%s", ip)
}

func getRedisToNumberKey(toNumber string) string {
	return fmt.Sprintf("APP_EMAIL_SMS_TO_%s", toNumber)
}

func getTotalEmailKey() string {
	return "APP_EMAIL_SMS_COUNTER"
}

func getReplyMessageKeyByID(id string) string {
	return "APP_EMAIL_SMS_REPLY_" + strings.Replace(data.Config.EmailConfig.EmailAddress, "@", fmt.Sprintf("+%s@", id), 1)
}

func getReplyMessageKeyByEmail(email string) string {
	return fmt.Sprintf("APP_EMAIL_SMS_REPLY_%s", email)
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
