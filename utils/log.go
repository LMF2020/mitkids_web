package utils

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	// 指定不同级别的log输出路径
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  "/go-docker/logs/info.log",
		logrus.ErrorLevel: "/go-docker/logs/error.log",
		logrus.DebugLevel: "/go-docker/logs/debug.log",
		logrus.PanicLevel: "/go-docker/logs/panic.log",
	}

	Log = logrus.New()
	// 输出文件行号
	Log.SetReportCaller(true)
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return Log
}
