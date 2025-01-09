package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

type mailProp struct {
	domain     string
	host       string
	port       int
	username   string
	password   string
	encryption string
	fromAddr   string
	fromName   string
}

type msgProp struct {
	from        string
	fromName    string
	to          string
	subject     string
	attachments []string
	data        any
	dataMap     map[string]any
}

func currentTime() string {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("now=%s\n", now)
	return now
}

func (appCtx applicationContext) sendSMTPMessage(msg msgProp) error {
	appCtx.logger.Printf("sendSMTPMessage msg=%+v", msg)

	if msg.from == "" {
		msg.from = appCtx.mailer.fromAddr
	}

	if msg.fromName == "" {
		msg.fromName = appCtx.mailer.fromName
	}

	data := map[string]any{
		"message": msg.data,
	}

	msg.dataMap = data

	appCtx.logger.Printf("msg.dataMap len=%d", len(msg.dataMap))

	for k, v := range msg.dataMap {
		appCtx.logger.Printf("%v=%v", k, v)
	}

	formattedMsg, err := appCtx.mailer.buildHTMLMessage(msg)

	if err != nil {
		appCtx.logger.Println("failed buildHTMLMessage", err)
		return err
	}

	plainMsg, err := appCtx.mailer.buildPlainTxtMessage(msg)

	if err != nil {
		appCtx.logger.Println("failed buildPlainTxtMessage", err)
		return err
	}

	appCtx.logger.Println("plainMsg=", plainMsg)
	appCtx.logger.Println("formattedMsg=", formattedMsg)

	mailSrv := mail.NewSMTPClient()

	mailSrv.Host = appCtx.mailer.host
	mailSrv.Port = appCtx.mailer.port
	mailSrv.Username = appCtx.mailer.username
	mailSrv.Password = appCtx.mailer.password
	mailSrv.Encryption = appCtx.mailer.getEncryption(appCtx.mailer.encryption)
	mailSrv.KeepAlive = false
	mailSrv.ConnectTimeout = 10 * time.Second
	mailSrv.SendTimeout = 10 * time.Second

	emsg := mail.NewMSG()

	emsg.SetFrom(msg.from).AddTo(msg.to).SetSubject(msg.subject + " " + currentTime())
	emsg.SetBody(mail.TextPlain, plainMsg)
	emsg.AddAlternative(mail.TextHTML, formattedMsg)

	if len(msg.attachments) > 0 {
		for _, x := range msg.attachments {
			emsg.AddAttachment(x)
		}
	}

	smtpClient, err := mailSrv.Connect()
	if err != nil {
		appCtx.logger.Printf("Failed to connected to mailSrv=%+v", mailSrv)
		appCtx.logger.Println("error=", err.Error())
		return err
	}

	appCtx.logger.Printf("mailSrv connected smtpClient=%+v", smtpClient)

	err = emsg.Send(smtpClient)

	if err != nil {
		return err
	}

	appCtx.logger.Printf("Sent successfully emsg=%+v", emsg)

	return nil
}

func (m mailProp) buildPlainTxtMessage(msg msgProp) (string, error) {
	templatePath := "./templates/mail-plain.gohtml" // mail-plain.gohtml

	t, err := template.New("email-plain").ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = t.ExecuteTemplate(&buf, "body", msg.dataMap); err != nil {
		return "", err
	}

	plain := buf.String()

	return plain, nil
}

func (m mailProp) buildHTMLMessage(msg msgProp) (string, error) {
	templatePath := "./templates/mail.gohtml" // ./templates/mail.gohtml"

	t, err := template.New("email").ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = t.ExecuteTemplate(&buf, "body", msg.dataMap); err != nil {
		return "", err
	}

	formatted := buf.String()

	formatted, err = m.inlineCSS(formatted)

	if err != nil {
		return "", err
	}

	return formatted, nil
}

func (m mailProp) inlineCSS(txt string) (string, error) {
	opts := premailer.Options{
		RemoveClasses: false, CssToAttributes: false, KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(txt, &opts)

	if err != nil {
		return "", err
	}

	html, err := prem.Transform()

	if err != nil {
		return "", err
	}

	return html, err
}

func (m mailProp) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
		// return mail.EncryptionSSLTLS
	default:
		return mail.EncryptionSTARTTLS
	}
}
