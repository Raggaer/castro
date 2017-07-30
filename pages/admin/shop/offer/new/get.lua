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

    data.category = db:singleQuery("SELECT id, name FROM castro_shop_categories WHERE id = ?", http.getValues.categoryId)

    if data.category == nil then
        http:redirect("/")
        return
    end

    http:render("newoffer.html", data)
end