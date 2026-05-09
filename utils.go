package main

import (
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	return logger
}

func ErrorBody(err string) gin.H {
	return gin.H{
		"error": err,
	}
}

func URLWithParams(base string, params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return base + "?" + values.Encode()
}

func Includes[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
