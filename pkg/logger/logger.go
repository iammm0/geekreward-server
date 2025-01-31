package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	InfoLogger  = logrus.New()
	ErrorLogger = logrus.New()
)

func InitLogger() {
	// 设置输出到标准输出
	InfoLogger.SetOutput(os.Stdout)

	// 设置日志级别为Info级别
	InfoLogger.SetLevel(logrus.InfoLevel)

	// 设置日志格式为JSON格式（可选，默认为文本格式）
	InfoLogger.SetFormatter(&logrus.JSONFormatter{})

	// 错误日志同理
	ErrorLogger.Out = os.Stderr
	ErrorLogger.SetLevel(logrus.ErrorLevel)
	ErrorLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 如果需要，可以添加文件或其他日志目标
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		InfoLogger.Out = file
	} else {
		InfoLogger.Info("Failed to log to file, using default stderr")
	}
}
