package handler

import (
	"blog_platform/cmd/server/internal/repository"	
	"net/http"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	Repo repository.ArticleRepository
}

func NewArticleHandler(repo repository.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{Repo: repo}
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