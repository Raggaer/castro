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

func ArticleSingle(query string, args []interface{}) (*Article, error) {
	article := Article{}

	// Get article using RAW query
	if err := database.DB.Table("articles").Where(query, args).Scan(&article).Error; err != nil {
		return nil, err
	}

	// Return all values
	return &article, nil
}

func SaveArticle(article *Article) error {
	return database.DB.Save(article).Error
}