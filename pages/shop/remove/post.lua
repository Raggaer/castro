function post()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local cart = session:get("shop-cart")

    if cart == nil then
        http:redirect("/")
        return
    end

    local i = 0

    for name, count in pairs(cart) do
        if name == http.postValues["offer"] then
            if count == 1 then
                cart[name] = nil
                session:set("shop-cart", cart)
                session:setFlash("success", "Offer removed from your cart")
                http:redirect("/subtopic/shop/view")
                return
            end
            cart[name] = count - 1
            session:set("shop-cart", cart)
            session:setFlash("success", "Offer removed from your cart")
            http:redirect("/subtopic/shop/view")
            return
        end
        i = i + 1
    end

    http:redirect("/")
end