package models

import (
	"time"
	"github.com/raggaer/castro/app/database"
)

// Article struct used to represent castro
// latest news
type Article struct {
	ID        int64
	Title     string
	Text      string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ArticleSingleWhere(where string, args ...interface{}) (*Article, error) {
	article := Article{}
	if err := database.DB.Table("articles").Where(where, args).Scan(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}