package logger

import (
	"fmt"
	"log"
	"os"
)

// SimpleLogger is the standar logger implementation
type SimpleLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	level       int
	LogSettings
}

// NewSimpleLogger creates a new instance of SimpleLogger.
func NewSimpleLogger() (sl *SimpleLogger) {
	sl = new(SimpleLogger)
	sl.level = LevelInfo

	sl.debugLogger = log.New(os.Stdout, "[DBG] ", log.LstdFlags)
	sl.infoLogger = log.New(os.Stdout, "[INF] ", log.LstdFlags)
	sl.errorLogger = log.New(os.Stdout, "[ERR] ", log.LstdFlags)
	return
}

// Debug prints the arguments to the debug logger.
func (sl *SimpleLogger) Debug(v ...interface{}) {
	if sl.GetLevel() >= LevelDebug {
		sl.debugLogger.Print(v...)
	}
}

// Debugf prints the arguments to the debug logger. Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Debugf(format string, v ...interface{}) {
	if sl.GetLevel() >= LevelDebug {
		sl.debugLogger.Printf(format, v...)
	}
}

// Info prints the arguments to the info logger.
func (sl *SimpleLogger) Info(v ...interface{}) {
	if sl.GetLevel() >= LevelInfo {
		sl.infoLogger.Print(v...)
	}
}

// Infof prints the arguments to the info logger. Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Infof(format string, v ...interface{}) {
	if sl.GetLevel() >= LevelInfo {
		sl.infoLogger.Printf(format, v...)
	}
}

// Error prints the arguments to the error logger.
func (sl *SimpleLogger) Error(v ...interface{}) {
	sl.errorLogger.Print(v...)
}

// Errorf prints the arguments to the error logger. Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Errorf(format string, v ...interface{}) {
	sl.errorLogger.Printf(format, v...)
}

// Fatal prints the arguments to the error logger, followed by a call to os.Exit(1).
func (sl *SimpleLogger) Fatal(v ...interface{}) {
	sl.errorLogger.Fatal(v...)
}

// Fatalf prints the arguments to the error logger, followed by a call to os.Exit(1).
// Arguments are handled like in fmt.Printf.
func (sl *SimpleLogger) Fatalf(format string, v ...interface{}) {
	sl.errorLogger.Fatalf(format, v...)
}

// Print prints the arguemtns to the info logger. It's good for standar logger compatibility
func (sl *SimpleLogger) Print(v ...interface{}) {
	sl.infoLogger.Print(v...)
}

// SetLevel sets the log level (0=ERROR, 1=INFO, 2=DEBUG)
func (sl *SimpleLogger) SetLevel(level int) {
	sl.level = level
}

// GetLevel returns the current log level
func (sl *SimpleLogger) GetLevel() int {
	return sl.level
}

// SetLogSettings sets the log settings.
func (sl *SimpleLogger) SetLogSettings(logSettings LogSettings) {
	sl.LogSettings = logSettings
}

// LogToSlack sends a message to the configured channel, if it's enabled
func (sl *SimpleLogger) LogToSlack(webHook, title, text string) {
	go func() {
		if sl.LogSettings == nil {
			sl.Fatalf("LogSettings have not been initialized.")
		}

		if sl.LogSettings.GetSlackEnabled() {
			if err := SendAlert(webHook, "", title, ColorGood, text, sl.LogSettings); err != nil {
				sl.Errorf("Found an error sending notification to Slack: %s", err)
			}
		}
	}()
}

// LogErrorToSlack sends a message formatted as error to the configured channel, if it's enabled
func (sl *SimpleLogger) LogErrorToSlack(webHook, title, text string) {
	go func() {
		if sl.LogSettings == nil {
			sl.Fatalf("LogSettings have not been initialized.")
		}

		if sl.LogSettings.GetSlackEnabled() {
			if err := SendAlert(webHook, "", title, ColorDanger, fmt.Sprintf("`%s`", text), sl.LogSettings); err != nil {
				sl.Errorf("Found an error sending notification to Slack: %s", err)
			}
		}
	}()
}
