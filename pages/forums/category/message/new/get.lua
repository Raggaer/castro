function get()
    if not session:isLogged() then
        http:redirect("/subtopic/forums")
        return
    end

    local data = {}
    local account = session:loggedAccount()

    data.characters = db:query("SELECT name, vocation, level FROM players WHERE account_id = ? ORDER BY id DESC", account.ID)
    data.info = db:singleQuery("SELECT id, title FROM castro_forum_post WHERE id = ?", http.getValues.id)

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    http:render("newmessage.html", data)
end