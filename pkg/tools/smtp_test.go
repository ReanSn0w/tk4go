package tools_test

import (
	"os"
	"testing"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

var (
	smtp *tools.SMTP

	smtpConfig = tools.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Login:    os.Getenv("SMTP_LOGIN"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Name:     os.Getenv("SMTP_NAME"),
		Email:    os.Getenv("SMTP_EMAIL"),
	}
)

func Test_SendPlainMail(t *testing.T) {
	err := getSmtp().SendTextEmail(
		"Дмитрий Папков",
		"papkovda@me.com",
		"Test mail",
		"Текст сообщения",
	)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SendHTMLMail(t *testing.T) {
	err := getSmtp().SendHTMLEmail(
		"Дмитрий Папков",
		"papkovda@me.com",
		"Test mail",
		[]byte("<html><body><h1>hello, world</h1></body></html>"),
	)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func getSmtp() *tools.SMTP {
	if smtp == nil {
		smtp = tools.NewSMTP(smtpConfig)
	}

	return smtp
}
