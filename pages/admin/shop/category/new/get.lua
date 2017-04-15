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

    data.success = session:getFlash("success")
    data.validationError = session:getFlash("validationError")

    http:render("shopcategory.html", data)
end