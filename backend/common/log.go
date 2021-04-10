package common

import (
	"go.uber.org/zap"
)

// Logger is used fo logging. Ensure SetUpLogger is called way before using it.
var Logger *zap.SugaredLogger

// SetUpLogger configures Logger
func SetUpLogger(isProd bool) {
	var l *zap.Logger
	if isProd {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	Logger = l.Sugar()
}
