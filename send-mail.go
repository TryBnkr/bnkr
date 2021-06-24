package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/MohammedAl-Mahdawi/bnkr/app/types"
	"github.com/MohammedAl-Mahdawi/bnkr/utils"

	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	// Run pool of 5 email workers which allows 5 emails to be sent at the same time
	for i := 0; i < 5; i++ {
		go func() {
			for {
				msg := <-app.MailChan
				sendMsg(msg)
			}
		}()
	}
}

func sendMsg(m types.MailData) {
	server := mail.NewSMTPClient()
	server.Host = utils.GetOptionValue("SMTP_HOST")
	port, err := strconv.Atoi(utils.GetOptionValue("SMTP_PORT"))
	if err != nil {
		return
	}
	server.Port = port
	server.Username = utils.GetOptionValue("SMTP_USERNAME")
	server.Password = utils.GetOptionValue("SMTP_PASSWORD")
	switch utils.GetOptionValue("SMTP_SECURITY") {
	case "none":
		server.Encryption = mail.EncryptionNone
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	case "tls":
		server.Encryption = mail.EncryptionTLS
	}
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To...).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	}
}
