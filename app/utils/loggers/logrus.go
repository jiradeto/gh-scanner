package loggers

import (
	"strings"

	"github.com/jiradeto/gh-scanner/app/environments"
	"github.com/sirupsen/logrus"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
// Default logrus set log level = InfoLevel
func NewLogger(mode string) *StandardLogger {
	var baseLogger = logrus.New()

	if environments.DevMode {
		baseLogger.SetLevel(logrus.DebugLevel)
	}

	var standardLogger = &StandardLogger{baseLogger}

	switch strings.ToLower(mode) {
	case "json":
		standardLogger.Formatter = &logrus.JSONFormatter{}
	case "text":
		standardLogger.Formatter = &logrus.TextFormatter{}
	default:
		standardLogger.Formatter = &logrus.JSONFormatter{}
	}

	return standardLogger
}
