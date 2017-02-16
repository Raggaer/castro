-- Block access for anyone who is not admin
if not session:isLogged() or not session:isAdmin() then
    http:redirect("/")
    return
end

local removeArticle = tonumber(http.postValues["id"])
if removeArticle then
	db:execute("DELETE from castro_articles WHERE id = ?", removeArticle)
	session:setFlash("success", "Article was successfully removed.")
end

http:redirect("/subtopic/admin/articles/list")