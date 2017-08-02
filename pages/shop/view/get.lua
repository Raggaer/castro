function get()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data.categories = db:query("SELECT id, name, description FROM castro_shop_categories ORDER BY id")

    for _, category in ipairs(data.categories) do
        category.offers = db:query("SELECT name, description, price FROM castro_shop_offers WHERE category_id = ?", category.id)
    end

    http:render("shopview.html", data)
end