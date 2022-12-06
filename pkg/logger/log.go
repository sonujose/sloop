package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// CorrelationHeader defines a default Correlation ID HTTP header.
const (
	LogTimestampFormat = "2006-01-02 15:04:05"
)

func NewLogger() *logrus.Logger {

	log = logrus.New()

	setLoggingVerbos(1)

	// Setting up format for logrus module
	log.SetFormatter(&logrus.TextFormatter{TimestampFormat: LogTimestampFormat,
		FullTimestamp: true, ForceColors: true, DisableLevelTruncation: true})

	log.SetLevel(getLowestLoggingLevel())
	log.SetOutput(logrus.StandardLogger().Out)

	return log

}

// Get the lowest logging level configuration
// If nothing specified or wrong value specified it will be default
func getLowestLoggingLevel() logrus.Level {

	lvl := os.Getenv("LOGGING_LEVEL")

	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}

	return ll
}

func setLoggingVerbos(level int) {
	var logLevel string
	if level > 1 {
		logLevel = "debug"
	} else {
		logLevel = "info"
	}
	os.Setenv("LOGGING_LEVEL", logLevel)
}
