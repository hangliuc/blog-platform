package main

import (
	"blog_platform/cmd/server/internal/handler"
	"blog_platform/cmd/server/internal/model"
	"blog_platform/cmd/server/internal/repository"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "root:GB_UXa>cX*h!K2@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("无法连接数据库: ", err)
	}

	db.AutoMigrate(&model.ArticleStat{})

	articleRepo := repository.NewArticleRepo(db)
	articleHandler := handler.NewArticleHandler(articleRepo)

	r := gin.Default()

	//  配置 CORS (允许跨域)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://hangops.top", "http://localhost:1313"}, // 允许的域名
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/article/visit", articleHandler.RecordArticleVisit)
	}

	log.Println("服务启动在 :8080")
	r.Run(":8080")
}