## 📬 Smaila – Simple Mail API

Smaila is a lightweight HTTP API for sending emails using an SMTP server. It supports HTML bodies, CC/BCC, and is ready to use with Docker. Swagger documentation is included out of the box.

## 🚀 Features

- Send emails with subject and body
- Support for HTML email content
- Multiple recipients (To, CC, BCC)
- Swagger (OpenAPI) documentation
- Dockerized for easy deployment

> ⚠️ Attachment support is currently limited – sending attachments is not functional at the moment.

## ⚙️ Setup

1. Create a .env file:

```env
HTTP_PORT=8080

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_FROM=example@gmail.com
SMTP_USERNAME=example@gmail.com
SMTP_PASSWORD=aaaabbbbccccdddd
```

2.  Pull the Docker image:

```bash
docker pull ghcr.io/cactusbros/smaila:latest # or specific version like v0.2.2
```

Or manually Build the image:

```bash
git clone https://github.com/CactusBros/smaila.git
cd smaila
docker build -t smaila .
```

3. Run the container:

```bash
docker run -p 8080:8080 --env-file .env smaila
```

## 📘 API Documentation

After running the container, access the interactive API docs at:
http://127.0.0.1:8080/swagger

## 📤 Example Usage

You can send a POST request to / using curl or tools like Postman:

```bash
curl -X POST http://localhost:8080/ \
  -F 'to=someone@example.com' \
  -F 'subject=Hello' \
  -F 'body=This is a test email' \
  -F 'is_html=false'
```

## 📝 Todo

- Support for sending attachments
- Add email open tracking
- Logging improvements
- Authentication or API key protection

## 🧑‍💻 Authors

Built with ❤️ by CactusBros
