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

    data.category = db:singleQuery("SELECT id, name, description, created_at FROM castro_shop_categories WHERE id = ?", http.getValues.id)

    if data.category == nil then
        http:redirect("/subtopic/admin/shop")
        return
    end

    data.validationError = session:getFlash("validationError")
    data.success = session:getFlash("success")
    data.list = db:query("SELECT id, name, price, offer_type FROM castro_shop_offers WHERE category_id = ?", http.getValues.id)

    http:render("shopoffers.html", data)
end