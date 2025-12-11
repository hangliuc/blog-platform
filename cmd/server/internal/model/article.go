package model

import "time"

// ArticleStat 对应数据库中的 article_stats 表
type ArticleStat struct {
	Path      string    `gorm:"primaryKey;type:varchar(191)"` 
	ViewCount int64     `gorm:"default:0"` // 浏览量
	CreatedAt time.Time // 创建时间 (自动维护)
	UpdatedAt time.Time // 更新时间 (自动维护)
}

type SiteStat struct {
	ID        uint  `gorm:"primaryKey"` 
	ViewCount int64 `gorm:"default:0"`
}