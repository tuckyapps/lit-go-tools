package api

import (
	"io/ioutil"
)

// API messages
const (
	MsgContactThanks         = "msg_contact_thanks"
	MsgContactThanks48Hours  = "msg_contact_thanks_48_hours"
	MsgThanksForYourFeedback = "thanks_for_your_feedback"
)

// Message constants
const (
	DefaultLanguage    = "es"
	TranslationsFolder = "lang/"
	SupportedLanguages = "es"
)

// structure to hold server messages in the different supported languages
var messages map[string]map[string]string

// GetMessage returns the text for msgID, in the specific language
func GetMessage(msgID, lang string) string {

	selectedLang := lang

	if lang == "" {
		selectedLang = DefaultLanguage
	}

	// check if language is supported; if not, use the default one
	_, exists := messages[selectedLang]
	if !exists {
		selectedLang = DefaultLanguage
	}

	holder := messages[selectedLang]
	return holder[msgID]
}

func getDirectoryFiles(path string) (fileNames []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() {
				fileNames = append(fileNames, f.Name())
			}
		}
	}

	return
}
