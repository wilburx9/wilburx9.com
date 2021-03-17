package common

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func SetUpLogger(isProd bool) {
	var l *zap.Logger
	if isProd {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	Logger = l.Sugar()
}

