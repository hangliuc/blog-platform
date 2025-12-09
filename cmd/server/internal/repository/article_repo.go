package repository

import (
	"blog_platform/cmd/server/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 定义接口，方便未来做单元测试或更换数据库
type ArticleRepository interface {
	IncreaseViewCount(path string) (int64, error)
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) ArticleRepository {
	return &articleRepo{db: db}
}

// IncreaseViewCount 原子性增加访问量并返回最新值
func (r *articleRepo) IncreaseViewCount(path string) (int64, error) {
	// 1. 执行 Upsert (存在则更新，不存在则插入)
	// SQL: INSERT INTO article_stats ... ON DUPLICATE KEY UPDATE view_count = view_count + 1
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "path"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"view_count": gorm.Expr("view_count + 1"),
		}),
	}).Create(&model.ArticleStat{Path: path, ViewCount: 1}).Error

	if err != nil {
		return 0, err
	}

	// 2. 查出最新的值返回
	var stat model.ArticleStat
	err = r.db.Select("view_count").Where("path = ?", path).First(&stat).Error
	return stat.ViewCount, err
}