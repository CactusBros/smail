package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var envPath = flag.String("env", ".env", "path to the environment variable file")

func main() {
	flag.Parse()
	if v := os.Getenv("CONFIG_PATH"); len(v) != 0 {
		*envPath = v
	}
	cfg := MustInitConfig()
	app := fiber.New()
	app.Post("/", GetMailHandler(cfg.SMTP))
	app.Listen(fmt.Sprintf(":%d", cfg.HTTP.Port))
}

type Config struct {
	HTTP HTTPConfig
	SMTP SMTPConfig
}

type HTTPConfig struct {
	Port int `env:"HTTP_PORT"`
}

type SMTPConfig struct {
	Host     string `env:"SMTP_HOST"`
	Port     int    `env:"SMTP_PORT"`
	From     string `env:"SMTP_FROM"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
}

func MustInitConfig() Config {
	err := godotenv.Load(*envPath)
	if err != nil {
		panic(err)
	}
	cfg, err := ReadConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

func ReadConfigFromEnv() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

type MailMessage struct {
	To          []string     // primary recipients
	Cc          []string     // carbon copy
	Bcc         []string     // blind carbon copy
	Subject     string       // email subject
	Body        string       // main message body (usually plain text or HTML)
	IsHTML      bool         // true if Body is HTML
	Attachments []Attachment // optional file attachments
}

type Attachment struct {
	Filename    string // name of the file
	ContentType string // MIME type, e.g. "application/pdf"
	Content     []byte // file data
}

func SendMail(cfg SMTPConfig, msg MailMessage) error {
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
	for _, a := range msg.Attachments {
		m.Attach(a.Filename)
	}

	// Send email
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return d.DialAndSend(m)
}

func GetMailHandler(cfg SMTPConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse basic fields from the form
		to := strings.Split(c.FormValue("to"), ",")
		cc := strings.Split(c.FormValue("cc"), ",")
		bcc := strings.Split(c.FormValue("bcc"), ",")
		subject := c.FormValue("subject")
		body := c.FormValue("body")
		isHTML := c.FormValue("is_html") == "true"

		// Handle attachments
		form, err := c.MultipartForm()
		if err != nil && err != http.ErrNotMultipart {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid multipart form")
		}

		var attachments []Attachment
		if form != nil {
			files := form.File["attachments"]
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Failed to open attachment")
				}
				defer file.Close()

				content, err := io.ReadAll(file)
				if err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "Failed to read attachment")
				}

				attachments = append(attachments, Attachment{
					Filename:    fileHeader.Filename,
					ContentType: fileHeader.Header.Get("Content-Type"),
					Content:     content,
				})
			}
		}

		// Compose the mail message
		msg := MailMessage{
			To:          filterEmpty(to),
			Cc:          filterEmpty(cc),
			Bcc:         filterEmpty(bcc),
			Subject:     subject,
			Body:        body,
			IsHTML:      isHTML,
			Attachments: attachments,
		}

		// Send
		if err := SendMail(cfg, msg); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func filterEmpty(input []string) []string {
	var result []string
	for _, s := range input {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
