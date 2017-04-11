function widget()

    if not session:isLogged() then
        widgets:render("admin.html", nil)
        return
    end

    local data = {}

    data.admin = session:loggedAccount().castro.Admin

    widgets:render("admin.html", data)
end