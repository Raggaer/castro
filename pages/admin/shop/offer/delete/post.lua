function post()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isAdmin() then
        http:redirect("/")
        return
    end

    local data = {}

    data.offer = db:singleQuery("SELECT id, category_id, price, description, name FROM castro_shop_offers WHERE id = ?", http.postValues.id)

    if data.offer == nil then
        http:redirect("/")
        return
    end

    db:execute("DELETE FROM castro_shop_offers WHERE id = ?", data.offer.id)

    session:setFlash("success", "Shop offer deleted")
    http:redirect("/subtopic/admin/shop/category?id=" .. data.offer.category_id)
end