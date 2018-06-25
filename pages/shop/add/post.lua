function post()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        session:setFlash("error", "You need to be logged in")
        http:redirect("/subtopic/shop/view")
        return
    end

    local offer = db:singleQuery("SELECT id, name, description, price FROM castro_shop_offers WHERE id = ?", http.postValues["offer"])

    if offer == nil then
        http:redirect("/")
        return
    end

    local shopCart = session:get("shop-cart")

    if shopCart == nil then
        local c = {}

        c[offer.name] = 1

        session:set("shop-cart", c)
        session:setFlash("success", "Offer added to your cart")

        http:redirect("/subtopic/shop/view")
        return
    end

    if shopCart[offer.name] then
        shopCart[offer.name] = shopCart[offer.name] + 1
    else
        shopCart[offer.name] = 1
    end

    session:set("shop-cart", shopCart)
    session:setFlash("success", "Offer added to your cart")

    http:redirect("/subtopic/shop/view")
end