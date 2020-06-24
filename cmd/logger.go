package main

import (
	log "github.com/sirupsen/logrus"
	"io"
)

func setupLogger() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
}

func setupFileLogger(file io.Writer) {
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
}
