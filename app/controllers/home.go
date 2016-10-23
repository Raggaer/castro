package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
)

// Home is the main aac homepage
func Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Load articles with the page 0
	articles, err := getArticles(0, 4)
	if err != nil {
		util.Logger.Info(err)
		return
	}

	// Render template
	util.Template.RenderTemplate(w, r, "home.html", map[string]interface{}{
		"articles": articles,
	})
}

// getArticles helper method to load articles with the given
// offset. To make pagination work
func getArticles(page, per int) ([]models.Article, error) {
	articles := []models.Article{}
	if err := database.DB.Table("articles").Offset(page).Limit(per).Scan(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}
