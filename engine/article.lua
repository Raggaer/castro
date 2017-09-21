require "bbcode"
require "extensionhooks"

function editArticle(editMode)
	-- Block access for anyone who is not admin
	if not session:isLogged() or not session:isAdmin() then
	    http:redirect("/")
	    return
	end

	local article = {
		id = tonumber(http.postValues["id"]),
		title = http.postValues["title"],
		text = http.postValues["text"],
		action = http.postValues["action"],
		editmode = editMode,
	}

	article.heading = (article.editmode == "edit" and "Edit article" or "New article")

	if not article.title or article.title:len() < 1 then
		article.validationError = "Title can not be empty."
		http:render("editarticle.html", article)
		return
	elseif not article.text or article.text:len() < 1 then
		article.validationError = "Text field can not be empty."
		http:render("editarticle.html", article)
		return
	end
	

	if article.action == "new" then
		db:execute("INSERT INTO castro_articles (title, text, created_at) VALUES (?, ?, NOW())", article.title, article.text)
		session:setFlash("success", "Article posted.")
		executeHook("onNewArticle", article, session:loggedAccount())
	elseif article.action == "edit" then
		db:execute("UPDATE castro_articles SET title = ?, text = ?, updated_at = NOW() WHERE id = ?", article.title, article.text, article.id)
		session:setFlash("success", "Article updated.")
        executeHook("onEditArticle", article, session:loggedAccount())
	end

	http:redirect("/subtopic/admin/articles/list")
end
