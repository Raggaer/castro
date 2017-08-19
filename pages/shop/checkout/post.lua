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

    local cartdata = {}
    local totalprice = 0

    for name, count in pairs(cart) do
        cartdata[name] = {}
        cartdata[name].offer = db:singleQuery("SELECT name, price FROM castro_shop_offers WHERE name = ?", name)

        if cartdata[name].offer == nil then
            http:redirect("/")
            return
        end

        totalprice = tonumber(cartdata[name].offer.price) * count

        cartdata[name].count = count
    end

    local account = session:loggedAccount()

    if account.castro.Points < totalprice then
        session:setFlash("error", "You need more points")
        http:redirect("/subtopic/shop/view")
        return
    end
end