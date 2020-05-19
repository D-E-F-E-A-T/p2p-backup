package smtp

import (
	"html/template"
	"log"
	"path"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type mail struct {
	*gomail.Dialer
	*template.Template
}

type MailOption struct {
	PreText   string
	AfterText string

	ButtonLink string
	ButtonText string

	Company     string
	Unsubscribe string
}

var Mail *mail

func init() {
	dialer := gomail.NewDialer(viper.GetString("smtp.host"),
		viper.GetInt("smtp.port"),
		viper.GetString("smtp.user"),
		viper.GetString("smtp.password"))

	tmpl, err := template.ParseGlob(path.Join(
		viper.GetString("smtp.templates"),
		"mail.html",
	))
	if err != nil {
		log.Panicf("Error with loading SMTP template: %v", err)
	}
	Mail = &mail{
		Dialer:   dialer,
		Template: tmpl,
	}
}

func (m mail) Send(to, subject, msg string) {
	message := gomail.NewMessage()
	message.SetHeader("From", viper.GetString("smtp.user"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", `[P2P] `+subject)
	message.SetBody("text/html; charset=UTF-8", msg)
	if err := m.DialAndSend(message); err != nil {
		log.Printf("Error sending email: %+v", err)
	}
}
