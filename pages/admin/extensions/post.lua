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

    local data = {}

    data.origin = app.Plugin.Origin

    try(
        function()
            data.list = json:unmarshal(http:get(data.origin .. "/rest/list/search?name=".. http.postValues.name .. "&p=0"))
        end,
        function()
            data.error = "Unable to retrieve extension list"
        end
    )

    http:render("extensions.html", data)
end
