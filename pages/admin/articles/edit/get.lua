function get()
-- Block access for anyone who is not admin
if not session:isLogged() or not session:isAdmin() then
    http:redirect("/")
    return
end

local id = tonumber(http.getValues.id)
if not id then
	session:setFlash("validationError", "No article specified for edit.")
	http:redirect("/subtopic/admin/articles/list")
	return
end

local data = db:singleQuery("SELECT id, title, text FROM castro_articles WHERE id = ?", id)
if not data then
	session:setFlash("validationError", "No article with specified id.")
	http:redirect("/subtopic/admin/articles/list")
	return
end

http:render("editarticle.html", data)
	end