package logger

import "net/http"

// Default log levels
const (
	LevelError = iota
	LevelInfo
	LevelDebug
)

// Logger presents a common interface for logger
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	SetLevel(level int)
	GetLevel() int
	LogToSlack(webHook, title, text string, logSettings LogSettings)
	LogErrorToSlack(webHook, title, text string, logSettings LogSettings)
}

//LogSettings interface to be implemented by other project settings.
type LogSettings interface {
	GetSlackEnabled() bool
	GetHTTPClient() (client *http.Client)
	GetAppName() string
}

var (
	logImpl Logger
)

// SetLogger sets the current logger
func SetLogger(l Logger) {
	logImpl = l
}

// GetLogger returns the current logger defined for the service
func GetLogger() Logger {
	if logImpl == nil {
		logImpl = NewSimpleLogger()
	}
	return logImpl
}

// SetLogLevel configures the application's log level
func SetLogLevel(level int) {
	if logImpl != nil {
		logImpl.SetLevel(level)
	}
}
