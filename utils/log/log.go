package log

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"mitkid_web/conf"
)

var Logger *logrus.Logger

func Init(c *conf.Config) {
	newLogger(c)
}
func newLogger(c *conf.Config) *logrus.Logger {
	if Logger != nil {
		return Logger
	}

	// 指定不同级别的log输出路径
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "/go-docker/logs/info.log",
		logrus.ErrorLevel: "/go-docker/logs/error.log",
		logrus.DebugLevel: "/go-docker/logs/debug.log",
		logrus.PanicLevel: "/go-docker/logs/panic.log",
	}

	Logger = logrus.New()
	// 输出文件行号
	Logger.SetReportCaller(true)
	Logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	level := c.Log.Level
	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	return Logger
}
