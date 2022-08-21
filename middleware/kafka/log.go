package kafka

import (
	"github.com/CoderBenson/go_study/tool/env"
	"github.com/CoderBenson/go_study/tool/log"
	"github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	options := []log.Option{
		log.WithConsole(),
		log.WithFilePath("build/log"),
	}
	if env.IsProduct() {
		options = append(options, log.WithFormater(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02-15:04:05.000",
		}))
	}
	logger = log.NewAppLogger("kafka", options...)
}
