package main

import (
	log "github.com/sirupsen/logrus"
	"io"
)

func setupLogger(file io.Writer) {
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
}
