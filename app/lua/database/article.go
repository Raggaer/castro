package database

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/models"
)

func articleMetaTable(article *models.Article, name string, L *lua.LState) *lua.LTable {
	// Create the metatable with the given name
	metatable := L.NewTypeMetatable(name)

	// Set article values
	L.SetField(metatable, "ID", lua.LNumber(article.ID))
	L.SetField(metatable, "Title", lua.LString(article.Title))
	L.SetField(metatable, "Text", lua.LString(article.Text))
	L.SetField(metatable, "CreatedAt", lua.LNumber(article.CreatedAt.Unix()))
	L.SetField(metatable, "UpdatedAt", lua.LNumber(article.CreatedAt.Unix()))

	// Return table
	return metatable
}

func ArticleSingle(L *lua.LState) int {
	article, err := models.ArticleSingleWhere("id = ?", 1)
	if err != nil {
		L.RaiseError("Cannot get article: %v", err)
	}
	L.Push(articleMetaTable(article, "article", L))
	return 1
}
