package log

import (
	"fmt"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func NewLogger(app string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02-15:04:05.000",
	})
	logger.SetReportCaller(true)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("app_%s.log", app),
		MaxSize:    10 * 1024,
		MaxBackups: 3,
		MaxAge:     31,
		Compress:   true,
	})
	return logger
}
