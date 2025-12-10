package handler

import (
	"blog_platform/cmd/server/internal/repository"	
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
)

type ArticleHandler struct {
	Repo repository.ArticleRepository
	StartTime time.Time
}

func NewArticleHandler(repo repository.ArticleRepository, startTime time.Time) *ArticleHandler {
	return &ArticleHandler{
		Repo:      repo,
		StartTime: startTime,
	}
}

// RecordArticleVisit 处理 POST /api/article/visit
func (h *ArticleHandler) RecordArticleVisit(c *gin.Context) {
	// 定义请求参数结构
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	// 解析 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供 path 参数"})
		return
	}

	// 调用数据库层
	count, err := h.Repo.IncreaseViewCount(req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"path":  req.Path,
		"views": count,
	})
}

// GetSiteInfo 处理 GET /api/site/info
func (h *ArticleHandler) GetSiteInfo(c *gin.Context) {
	// 从数据库查统计
	totalViews, totalArticles, err := h.Repo.GetSiteStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	// 计算运行时间 (当前时间 - 启动时间)
	uptime := time.Since(h.StartTime)

	c.JSON(http.StatusOK, gin.H{
		"total_views":    totalViews,
		"total_articles": totalArticles,
		"uptime_seconds": int64(uptime.Seconds()), // 返回秒数，前端处理成 "xx天xx小时" 更灵活
		"start_time":     h.StartTime.Format(time.RFC3339),
	})
}