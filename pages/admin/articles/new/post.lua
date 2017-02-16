-- Block access for anyone who is not admin
if not session:isLogged() or not session:isAdmin() then
    http:redirect("/")
    return
end

local article = {
	title = http.postValues["title"],
	text = http.postValues["text"],
	created = os.time(),
}

if not article.title or article.title:len() < 1 then
	article.validationError = "Title can not be empty."
	http:render("newarticle.html", article)
	return
elseif not article.text or article.text:len() < 1 then
	article.validationError = "Text field can not be empty."
	http:render("newarticle.html", article)
	return
end

db:execute("INSERT INTO castro_articles (title, text, created_at) VALUES (?, ?, NOW())", article.title, article.text)

session:setFlash("success", "Article posted.")

http:redirect("/subtopic/admin/articles/list")