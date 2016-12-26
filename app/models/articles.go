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

	// Get article using WHERE query
	if err := database.DB.Table("articles").Raw(query, args).Scan(&article).Error; err != nil {
		return nil, err
	}

	// Return all values
	return &article, nil
}

func ArticleMultiple(query string, args []interface{}) ([]Article, error) {
	articleList := []Article{}

	// Get article list using WHERE query
	if err := database.DB.Table("articles").Raw(query, args).Scan(&articleList).Error; err != nil {
		return nil, err
	}

	// Return all values
	return articleList, nil
}

func SaveArticle(article *Article) error {
	return database.DB.Model(article).Updates(article).Error
}