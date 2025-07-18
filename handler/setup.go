package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/CactusBros/smaila/config"
	"github.com/CactusBros/smaila/internal/mail"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/CactusBros/smaila/docs"
)

// Run initial routes and serve HTTP requests.
//
//	@title			Mail Sending Service API
//	@version		v0.1.0
//	@description	A simple API to send emails with optional attachments
//	@host			localhost
//	@BasePath		/
func Run(cfg config.Config) error {
	app := fiber.New()
	app.Post("/", newMailHandler(cfg.SMTP))
	app.Get("/swagger/*", swagger.HandlerDefault)

	return app.Listen(fmt.Sprintf(":%d", cfg.HTTP.Port))
}

// NewMailHandler handles sending emails
//
//	@Summary	Send an email
//	@Tags		mail
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		to			formData	string	true	"Recipient(s), comma separated"
//	@Param		cc			formData	string	false	"CC recipient(s), comma separated"
//	@Param		bcc			formData	string	false	"BCC recipient(s), comma separated"
//	@Param		subject		formData	string	true	"Email subject"
//	@Param		body		formData	string	true	"Email body"
//	@Param		is_html		formData	string	false	"Set to 'true' if body is HTML"
//	@Param		attachments	formData	file	false	"Attachments (can upload multiple)"
//	@Success	200			{string}	string	"OK"
//	@Failure	400			{string}	string	"Bad Request"
//	@Failure	500			{string}	string	"Internal Server Error"
//	@Router		/ [post]
func newMailHandler(cfg config.SMTPConfig) fiber.Handler {
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
