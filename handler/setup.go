package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/CactusBros/smaila/config"
	"github.com/CactusBros/smaila/internal/mail"
	"github.com/gofiber/fiber/v2"
)

func Run(cfg config.Config) error {
	app := fiber.New()
	app.Post("/", getMailHandler(cfg.SMTP))
	return app.Listen(fmt.Sprintf(":%d", cfg.HTTP.Port))
}

func getMailHandler(cfg config.SMTPConfig) fiber.Handler {
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

		var attachments []mail.Attachment
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

				attachments = append(attachments, mail.Attachment{
					Filename:    fileHeader.Filename,
					ContentType: fileHeader.Header.Get("Content-Type"),
					Content:     content,
				})
			}
		}

		// Compose the mail message
		msg := mail.Message{
			To:          filterEmpty(to),
			Cc:          filterEmpty(cc),
			Bcc:         filterEmpty(bcc),
			Subject:     subject,
			Body:        body,
			IsHTML:      isHTML,
			Attachments: attachments,
		}

		// Send
		if err := mail.Send(cfg, msg); err != nil {
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
