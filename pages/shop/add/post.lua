function post()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    local offer = db:singleQuery("SELECT id, name, description, price FROM castro_shop_offers WHERE id = ?", http.postValues["offer"])

    if offer == nil then
        http:redirect("/")
        return
    end

    local shopCart = session:get("shop-cart")

    if shopCart == nil then
        local cart = {}

        table.insert(cart, offer.id)

        session:set("shop-cart", cart)
        session:setFlash("success", "Offer added to your cart")

        http:redirect("/subtopic/shop/view")
        return
    end

    table.insert(shopCart, offer.id)

    session:set("shop-cart", cart)
    session:setFlash("success", "Offer added to your cart")

    http:redirect("/subtopic/shop/view")
end