package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"douCSAce/app/author"
	"douCSAce/app/confInstance"
	"douCSAce/app/confSeries"
	"douCSAce/app/journal"
	"douCSAce/app/paper"
	"douCSAce/pkg"

	_ "douCSAce/docs"
)

func initial(confFilepath string, logLevel string) {
	pkg.Conf = pkg.NewConfig(confFilepath)
	pkg.Log = pkg.NewLog(pkg.Conf.LogPath, logLevel)
	pkg.DB = pkg.NewDB(pkg.Conf.ArangoDB.Url, pkg.Conf.ArangoDB.Username, pkg.Conf.ArangoDB.Passwd,
		pkg.Conf.ArangoDB.Database, pkg.Conf.ArangoDB.ColNames)
}

func main() {
	// 初始化
	confFilepath := "conf.yaml"
	logLevel := pkg.DebugLevel
	initial(confFilepath, logLevel)

	// 初始化路由
	router := gin.Default()

	// 注册 swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 配置路由
	v1 := router.Group("/api/v1")
	{
		authorRouter := v1.Group("/author")
		{
			authorRouter.GET("/count", author.Count)
		}
		confInsRouter := v1.Group("/confInstance")
		{
			confInsRouter.GET("/count", confInstance.Count)
		}
		confSeriesRouter := v1.Group("/confSeries")
		{
			confSeriesRouter.GET("/count", confSeries.Count)
		}
		// fieldRouter := v1.Group("/field")
		journalRouter := v1.Group("/journal")
		{
			journalRouter.GET("/count", journal.Count)
		}
		paperRouter := v1.Group("/paper")
		{
			paperRouter.GET("/count", paper.Count)
		}
	}

	// 运行
	router.Run(pkg.Conf.Addr)
}
