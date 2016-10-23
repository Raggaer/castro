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
	articles, err := getArticles(0, 4)
	if err != nil {
		util.Logger.Info(err)
		return
	}
	util.Template.RenderTemplate(w, r, "home.html", map[string]interface{}{
		"articles": articles,
	})
}

func getArticles(page, per int) ([]models.Article, error) {
	articles := []models.Article{}
	if err := database.DB.Table("articles").Offset(page).Limit(per).Scan(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}
