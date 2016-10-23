package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
)

// Home is the main aac homepage
func Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	articles, total, err := getArticles(0, 4)
	if err != nil {
		util.Logger.Info(err)
		return
	}
	util.Logger.Info(articles, total)
	util.Template.RenderTemplate(w, r, "home.html", nil)
}

func getArticles(page, per int) ([]models.Article, int, error) {
	articles := []models.Article{}

	return articles, 0, nil
}
