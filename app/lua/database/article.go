package database

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/models"
	"strings"
	"strconv"
	"time"
)

var articleMethods = map[string]lua.LGFunction{
	"save": articleSave,
}

func articleMetaTable(article *models.Article, name string, L *lua.LState) *lua.LTable {
	// Create the metatable with the given name
	metatable := L.NewTypeMetatable(name)

	// Set article values
	L.SetField(metatable, "ID", lua.LNumber(article.ID))
	L.SetField(metatable, "Title", lua.LString(article.Title))
	L.SetField(metatable, "Text", lua.LString(article.Text))
	L.SetField(metatable, "CreatedAt", lua.LNumber(article.CreatedAt.Unix()))
	L.SetField(metatable, "UpdatedAt", lua.LNumber(article.UpdatedAt.Unix()))

	// Create user data
	articleData := L.NewUserData()
	articleData.Value = article

	// Set user data
	L.SetField(metatable, "__article", articleData)

	// Set metatable functions
	L.SetFuncs(metatable, articleMethods)

	// Return table
	return metatable
}

func metatableArticle(L *lua.LState, metatable *lua.LTable) (models.Article, error) {
	article := models.Article{}

	// Try to parse ID field as int64
	id, err := strconv.ParseInt(L.GetField(metatable, "ID").String(), 10, 64)

	if err != nil {
		return article, err
	}

	// Set article ID
	article.ID = id

	// Set article title and text
	article.Title = L.GetField(metatable, "Title").String()
	article.Text =  L.GetField(metatable, "Text").String()

	// Try to parse CreatedAt field as int64
	createdUnix, err := strconv.ParseInt(L.GetField(metatable, "CreatedAt").String(), 10, 64)

	if err != nil {
		return article, err
	}

	// Convert unix timestamp to time.Time
	article.CreatedAt = time.Unix(createdUnix, 0)

	// Try to parse UpdatedAt field as int64
	updatedUnix, err := strconv.ParseInt(L.GetField(metatable, "UpdatedAt").String(), 10, 64)

	if err != nil {
		return article, err
	}

	// Convert unix timestamp to time.Time
	article.UpdatedAt = time.Unix(updatedUnix, 0)

	return article, nil
}

func articleSave(L *lua.LState) int {
	// Get metatable
	metatable := L.Get(1)

	// Check if its a valid metatable
	if metatable.Type() != lua.LTTable {

		// Raise error
		L.RaiseError("Cannot get article metatable")

		return 0
	}

	article, err := metatableArticle(L, metatable.(*lua.LTable))

	// Check for errors
	if err != nil {

		// Raise error
		L.RaiseError("Cannot save article: %v", err)

		return 0
	}


	// Try to save the article pointer
	if err := models.SaveArticle(&article); err != nil {

		// Raise error
		L.RaiseError("Cannot save article: %v", err)

		return 0
	}

	return 0
}

//
func ArticleSingle(L *lua.LState) int {
	// Get query
	query := L.Get(2)

	// Check if query is valid
	if query.Type() != lua.LTString {

		// Raise error
		L.RaiseError("Cannot get article: missing QUERY")
		return 0
	}

	// Count number of params
	n := strings.Count(query.String(), "?")

	args := []interface{}{}

	// Get all arguments matching the number of params
	for i := 0; i < n; i++ {

		// Append argument to slice
		args = append(args, L.Get(2 + n).String())
	}

	// Get article from database
	article, err := models.ArticleSingle(query.String(), args)
	if err != nil {

		// Raise error if needed
		L.RaiseError("Cannot get single article: %v", err)
	}

	// Push article metatable to stack
	L.Push(articleMetaTable(article, "article", L))

	return 1
}
