function get()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    if not app.Plugin.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data.list = json:unmarshal(http:get(app.Plugin.Origin .. "/test"))

    http:render("extensions.html", data)
end