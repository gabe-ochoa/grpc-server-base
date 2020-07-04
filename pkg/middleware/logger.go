package middleware

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LogFormatText = "text"
	LogFormatJson = "json"
)

func SetupLogger(format, level string) {
	// Configure logging
	if format == LogFormatJson {
		logrus.SetFormatter(JSONFormatter())
	} else {
		logrus.SetFormatter(TextFormatter())
	}

	// Only log to stdout for now
	logrus.SetOutput(os.Stdout)

	// Parse and set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Log level string not supported. Defaulting to debug logging.")
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)

	logrus.WithFields(logrus.Fields{
		"level":  level,
		"format": format,
	}).Info("Logger configured")
}

func JSONFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}
}

func TextFormatter() logrus.Formatter {
	return &logrus.TextFormatter{
		FullTimestamp: true,
		// log in a more human friendly format
		TimestampFormat: time.RFC1123,
	}
}
