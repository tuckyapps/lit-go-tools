package logger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// Slack colors for messages
const (
	ColorGood    = "good"
	ColorDanger  = "danger"
	ColorWarning = "warning"
)

// SendAlert sends a notification to the specified slack channel
func SendAlert(channel, username, title, color, text string, settings Settings) (err error) {

	if channel != "" {

		// control parameters are valid
		if username == "" {
			username = settings.GetAppName()
		}
		if color == "" {
			color = ColorGood
		}

		template := `
		{
			"username": "$USERNAME",
			"attachments": [
				{
					"title": "$TITLE",
					"color": "$COLOR",
					"text": "$TEXT"
				}
			]
		}
		`

		// replace custom data
		msg := strings.Replace(template, "$USERNAME", username, 1)
		msg = strings.Replace(msg, "$TITLE", title, 1)
		msg = strings.Replace(msg, "$COLOR", color, 1)
		msg = strings.Replace(msg, "$TEXT", text, 1)

		// send message using HTTP
		agent := gorequest.New()
		agent.Client = settings.GetHTTPClient()

		response, _, errPost := agent.Post(channel).Send(msg).End()

		if response != nil {
			defer response.Body.Close()
		}

		if errPost == nil {
			if response.StatusCode != http.StatusOK {
				err = fmt.Errorf("Slack returned status %s", response.Status)
			}
		} else {
			err = fmt.Errorf("Error sending message to Slack: %s", errPost[0])
		}
	} else {
		err = fmt.Errorf("Invalid channel")
	}

	return
}
