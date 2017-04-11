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

    local c = {}

    c.method = "POST"
    c.url = app.Plugin.Origin .. "/plugin/view/" .. http.getValues.id .. "/download"
    c.timeout = "5s"
    c.authentication = {}
    c.authentication.username = app.Plugin.Username
    c.authentication.password = app.Plugin.Password

    local resp, _, code = http:curl(c)

    if code == 500 then
        session:setFlash("Error", json:unmarshal(resp).Message)
        http:redirect("/subtopic/admin/extensions/view?id=" .. http.getValues.id)
        return
    end

    http:setHeader("Content-Type", "application/zip")
    http:setHeader("Content-Disposition", "attachment; filename=plugin.zip")
    http:write(resp)
end