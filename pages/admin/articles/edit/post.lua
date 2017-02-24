function post()
-- Block access for anyone who is not admin
if not session:isLogged() or not session:isAdmin() then
    http:redirect("/")
    return
end

local article = {
	id = tonumber(http.postValues["id"]),
	title = http.postValues["title"],
	text = http.postValues["text"],
	created = os.time(),
}

if not article.title or article.title:len() < 1 then
	session:setFlash("validationError", "Title can not be empty.")
	http:redirect("/subtopic/admin/articles/list")
	return
elseif not article.text or article.text:len() < 1 then
	session:setFlash("validationError", "Text field can not be empty.")
	http:redirect("/subtopic/admin/articles/list")
	return
end

db:execute("UPDATE castro_articles SET title = ?, text = ?, updated_at = NOW() WHERE id = ?", article.title, article.text, article.id)

session:setFlash("success", "Article updated.")

http:redirect("/subtopic/admin/articles/list")
	end