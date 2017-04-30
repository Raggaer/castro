function get()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isAdmin() then
        http:redirect("/")
        return
    end

    local data = {}

    data.validationError = session:getFlash("validationError")
    data.current = os.date("%Y-%m-%d")

    http:render("newdiscount.html", data)
end