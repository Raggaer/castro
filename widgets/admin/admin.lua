function widget()

    if not session:isLogged() then
        widgets:render("admin.html", nil)
        return
    end

    local data = {}

    data.admin = session:loggedAccount().castro.Admin
    data.pluginEnabled = app.Plugin.Enabled

    widgets:render("admin.html", data)
end