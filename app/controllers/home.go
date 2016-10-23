package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
)

// Home is the main aac homepage
func Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Articles per page
	perPage := 4

	// Get page param
	page := ps.ByName("page")

	if page == "" {
		// Load articles with the page 0
		articles, max, err := getArticles(0, perPage)
		if err != nil {
			util.Logger.Error(err)
			return
		}

		// Render template
		util.Template.RenderTemplate(w, r, "home.html", map[string]interface{}{
			"articles": articles,
			"page":     0,
			"max":      max / perPage,
		})
		return
	}

	// Convert page to int
	pageNumber, err := strconv.Atoi(page)

	// Check if pageNumber can be a valid page
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	// Get article list for the given page
	articles, max, err := getArticles(pageNumber, perPage)
	if err != nil {
		util.Logger.Error(err)
		return
	}

	// Check if there are any articles for this page
	if len(articles) <= 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	// Render template
	util.Template.RenderTemplate(w, r, "home.html", map[string]interface{}{
		"articles": articles,
		"page":     pageNumber,
		"max":      max / perPage,
	})
}

// getArticles helper method to load articles with the given
// offset. To make pagination work
func getArticles(page, per int) ([]models.Article, int, error) {
	articles := []models.Article{}
	maxArticles := 0
	if err := database.DB.Table("articles").Count(&maxArticles).Offset(page).Limit(per).Order("id DESC").Scan(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, maxArticles, nil
}
