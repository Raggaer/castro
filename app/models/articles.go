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

// GetArticleWhere gets a single article from
// the database by the given where query
func GetAricleWhere(where string, args []interface{}) (Article, error) {
	article := Article{}
	db := database.DB.Table("articles").Where(where, args).Scan(&article)
	return article, db.Error
}