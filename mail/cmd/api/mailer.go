package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain, Host                   string
	Username, Password, Encryption string
	FromAddress, FromName          string
	Post                           int
}

type Message struct {
	From, FromName, To string
	Subject            string
	Attachments        []string
	Data               any
	DataMap            map[string]any
}

/*
We can use just build in libraries to send email but we'll user some 3rd party libraries
	1. github.com/vanng822/go-premailer/premailer
	2. github.com/xhit/go-simple-mail/v2
*/

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Post
	server.Username = m.Username
	server.Password = m.Password
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Encryption = m.getEncryption(m.Encryption)

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// adding attachments
	if len(msg.Attachments) > 0 {
		for _, asset := range msg.Attachments {
			email.AddAttachment(asset)
		}
	}

	if err := email.Send(smtpClient); err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	emailTemplateFilePath := "./templates/mail.html.gohtml"
	t, err := template.New("email-html").ParseFiles(emailTemplateFilePath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)

	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	emailTemplateFilePath := "./templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(emailTemplateFilePath)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) inlineCSS(fm string) (string, error) {

	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(fm, &options)

	if err != nil {
		return "", err
	}

	html, err := prem.Transform()

	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(en string) mail.Encryption {
	switch en {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
