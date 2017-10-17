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

    local discount = db:singleQuery("SELECT id, valid_till, discount, uses, unlimited FROM castro_shop_discounts WHERE code = ?", http.postValues.discount)

    if discount ~= nil then
        if os.time() < tonumber(discount.valid_till) then
            if discount.unlimited or tonumber(discount.uses) > 0 then
                totalprice = totalprice - ((tonumber(discount.discount) * totalprice) / 100)
                db:execute("UPDATE castro_shop_discounts SET uses = uses - 1 WHERE id = ?", discount.id)
            end
        end
    end

    local account = session:loggedAccount()

    if account.castro.Points < totalprice then
        session:setFlash("error", "You need more points")
        http:redirect("/subtopic/shop/view")
        return
    end

    db:execute("UPDATE castro_accounts SET points = points - ? WHERE account_id = ?", totalprice, account.ID)
    session:set("shop-cart", {})
    session:setFlash("success", "You paid " .. totalprice .. " for all your cart items")
    http:redirect("/subtopic/shop/view")
end