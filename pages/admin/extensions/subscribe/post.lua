function post()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    if not app.Plugin.Enabled then
        http:redirect("/")
        return
    end

    local info = json:unmarshal(http:get(app.Plugin.Origin .. "/plugin/view/" .. http.getValues.id))

    if info.Error then
        session:setFlash("Error", data.info.Message)
        http:redirect("/subtopic/admin/extensions")
        return
    end

    if db:singleQuery("SELECT 1 FROM castro_extension_subscribe WHERE plugin_id = ?", info.ID) ~= nil then
        db:execute("DELETE FROM castro_extension_subscribe WHERE plugin_id = ?", info.ID)
        session:setFlash("Success", "Unsubscribed from " .. info.Name)
        http:redirect("/subtopic/admin/extensions/view?id=" .. info.ID)
        return
    end

    local last = time:parseDate(info.UpdatedAt, "2006-01-02T15:04:05Z")

    db:execute("INSERT INTO castro_extension_subscribe (plugin_id, name, updated_at) VALUES (?, ?, ?)", info.ID, info.Name, last)
    session:setFlash("Success", "Subscribed to " .. info.Name)
    http:redirect("/subtopic/admin/extensions/view?id=" .. info.ID)
end