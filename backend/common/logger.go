package common

import (
	"github.com/getsentry/sentry-go"
	"log"
)

// LogMsg calls sentry.CaptureMessage in release anf log.Println in debug
func LogMsg(msg string) {
	if Config.isRelease() {
		sentry.CaptureMessage(msg)
	} else {
		log.Println(msg)
	}
}

// LogError calls sentry.CaptureException in release anf log.Println in debug
func LogError(err error) {
	if Config.isRelease() {
		sentry.CaptureException(err)
	} else {
		log.Println(err)
	}
}
