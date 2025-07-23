package mail

import (
	"github.com/CactusBros/smaila/config"
	"gopkg.in/gomail.v2"
)

type Message struct {
	To          []string // primary recipients
	Cc          []string // carbon copy
	Bcc         []string // blind carbon copy
	Subject     string   // email subject
	Body        string   // main message body (usually plain text or HTML)
	IsHTML      bool     // true if Body is HTML
	Attachments []string // optional path to file attachments
}

type Attachment struct {
	Filename    string // name of the file
	ContentType string // MIME type, e.g. "application/pdf"
	Content     []byte // file data
}

func Send(cfg config.SMTPConfig, msg Message) error {
	m := gomail.NewMessage()

	m.SetHeader("From", cfg.From)
	m.SetHeader("To", msg.To...)
	if len(msg.Cc) > 0 {
		m.SetHeader("Cc", msg.Cc...)
	}
	if len(msg.Bcc) > 0 {
		m.SetHeader("Bcc", msg.Bcc...)
	}
	m.SetHeader("Subject", msg.Subject)

	// Set body
	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	// Attach files
	for _, att := range msg.Attachments {
		m.Attach(att)
	}

	// Send email
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return d.DialAndSend(m)
}
