package middleware

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"overseas-official-website/base"
	"strconv"
	"time"
)

func Logger() gin.HandlerFunc {
	logClient := logrus.New()

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logrus.Fatalf("open file err: %v", err)
	}

	logClient.Out = src
	logClient.SetLevel(logrus.Level(base.GConf.Server.LogLevel))
	apiLogPath := base.GConf.Server.LogRoot + "overseas_official_website"

	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d.log",
		//rotatelogs.WithLinkName(apiLogPath),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(30*24*time.Hour),    // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.InfoLevel:  logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
	})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		start := time.Now()
		c.Set("log", logClient)
		// 设置多语言, 默认为中文
		defaultLanguage := 1
		language := c.GetHeader("language")
		if language != "" {
			defaultLanguage, err = strconv.Atoi(language)
			if err != nil {
				// 转换出错 默认为中文
				defaultLanguage = 1
			}
		}
		c.Set("language", defaultLanguage)

		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		logClient.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
