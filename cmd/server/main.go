package main

import (
	"blog_platform/cmd/server/internal/handler"
	"blog_platform/cmd/server/internal/model"
	"blog_platform/cmd/server/internal/repository"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN 环境变量未设置")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接数据库: ", err)
	}
	siteFoundingDate := time.Date(2025, 12, 10, 0, 0, 0, 0, time.Local)

	db.AutoMigrate(&model.ArticleStat{}, &model.SiteStat{})

	articleRepo := repository.NewArticleRepo(db)
	articleHandler := handler.NewArticleHandler(articleRepo, siteFoundingDate)

	r := gin.Default()

	//  配置 CORS (允许跨域)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://hangops.top", "http://localhost:1313"}, 
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/article/visit", articleHandler.RecordArticleVisit)
		api.POST("/site/visit", articleHandler.RecordSiteVisit)
		api.GET("/site/info", articleHandler.GetSiteInfo)
	}

	log.Println("服务启动在 :8080")
	r.Run(":8080")
}