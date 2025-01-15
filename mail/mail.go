package mail

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gopkg.in/gomail.v2"
)

/*
@params verificationCode string, to string
@returns bool, error
if the verification code is sent successfully then return true, nil else return false, error message.
*/
func SendMail(verificationCode string, to string, username string) error {
	_, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-z]+\.[a-zA-Z]{2,}$`, to)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verification code for aruka feedback")
	m.SetBody("text/html", MailTemplate(verificationCode, username))

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid smtp port")
	}
	d := gomail.NewDialer(smtpServer, port, smtpUser, smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email")
	}
	return nil
}
