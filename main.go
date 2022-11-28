package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"overseas-official-website/api/refresh"
	"overseas-official-website/base"
	"overseas-official-website/middleware"
	"overseas-official-website/route"
)

func main() {
	// 初始化
	base.InitAppService()
	//启动定时任务
	go refresh.AppKeyRefresh()

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	engine.Use(gin.Recovery(), middleware.Logger())

	route.InitRouters(engine)

	logrus.Fatal(engine.Run(base.GConf.Server.Host))
}
