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

	c.JSON(http.StatusOK, gin.H{
		"path":  req.Path,
		"views": count,
	})
}

func (h *ArticleHandler) RecordSiteVisit(c *gin.Context) {
	// 调用 Repo 增加全站计数
	err := h.Repo.IncreaseSiteVisit()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

// GetSiteInfo 处理 GET /api/site/info
func (h *ArticleHandler) GetSiteInfo(c *gin.Context) {
	// 从数据库查统计
	totalViews, siteTotalVisits, err := h.Repo.GetSiteStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	uptime := time.Since(h.StartTime)

	c.JSON(http.StatusOK, gin.H{
		"total_views":    totalViews,
		"site_total_visits": siteTotalVisits,
		"uptime_seconds": int64(uptime.Seconds()), 
		"start_time":     h.StartTime.Format(time.RFC3339),
	})
}