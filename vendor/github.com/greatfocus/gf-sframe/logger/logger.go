package logger

import (
	"log"
	"os"

	logrus "github.com/sirupsen/logrus"
)

// NewLogger create instance of the logger types
func NewLogger(service string) *logrus.Logger {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "@caller",
		},
	})
	logger.WithFields(
		logrus.Fields{
			"service": service,
		})
	logger.SetLevel(logrus.TraceLevel)
	log.SetOutput(os.Stdout)
	return logger
}
