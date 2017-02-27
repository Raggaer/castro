function get()
	-- Block access for anyone who is not admin
	if not session:isLogged() or not session:isAdmin() then
		http:redirect("/")
		return
	end

	local data, edit = {}, tonumber(http.getValues.id)
	if edit ~= nil then
		data = db:singleQuery("SELECT id, title, text FROM castro_articles WHERE id = ?", edit)
		if not data then
			session:setFlash("validationError", "No article with specified id.")
			http:redirect("/subtopic/admin/articles/list")
			return
		end
	end

	data.editmode = "edit"
	data.heading = "Edit article"

	http:render("editarticle.html", data)
end