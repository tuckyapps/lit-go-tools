package email

import (
	"bytes"
	"text/template"

	"gopkg.in/gomail.v2"
)

// Email options.
type Email struct {
	// From is the source email.
	FromAddress string

	// FromName is the friendly name of the source.
	FromName string

	// To is a set of destination emails.
	To []string

	// ReplyTo is a set of reply to emails.
	ReplyTo []string

	// Subject is the email subject text.
	Subject string

	// Text is the plain text representation of the body.
	Text string

	// HTML is the HTML representation of the body.
	HTML string
}

// SMTPCredentials is used to hold the smtp credentials.
type SMTPCredentials struct {
	Host     string
	Port     int
	User     string
	Password string
}

// SendEmail an email.
func SendEmail(e *Email, c *SMTPCredentials) error {

	// Create a new message.
	m := gomail.NewMessage()

	m.SetBody("text/html", e.HTML)

	if e.HTML == "" {
		m.AddAlternative("text/plain", e.Text)
	}

	// Construct the message headers, including a Configuration Set and a Tag.
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(e.FromAddress, e.FromName)},
		"To":      e.To,
		"Subject": {e.Subject},
	})

	// Send the email.
	d := gomail.NewPlainDialer(c.Host, c.Port, c.User, c.Password)

	err := d.DialAndSend(m)

	return err
}

// ParseHTMLTemplate is used to parse an html file and include the data on it.
func ParseHTMLTemplate(templateFileName string, data interface{}) (body string, err error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body = buf.String()
	return body, nil
}
