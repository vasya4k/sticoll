package main

import (
	"github.com/sirupsen/logrus"
)

func logFatal(topic string, msg string, err error) {
	logrus.WithFields(logrus.Fields{
		"topic": topic,
		"event": msg,
	}).Fatal(err)
}
