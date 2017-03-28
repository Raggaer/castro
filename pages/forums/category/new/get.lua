function get()
    if not session:isLogged() then
        http:redirect("/subtopic/forums")
        return
    end

    local data = {}

    data["validationError"] = session:getFlash("validationError")
    data["success"] = session:getFlash("success")
    data.logged = session:isLogged()
    data.info = db:singleQuery("SELECT id, title, description FROM castro_forum_category WHERE id = ?", http.getValues.id)
    data.characters = db:query("SELECT name FROM players WHERE account_id = ?", session:loggedAccount().ID)

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    http:render("newcategory.html", data)
end