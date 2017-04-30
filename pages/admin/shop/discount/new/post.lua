function post()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isAdmin() then
        http:redirect("/")
        return
    end

    local check = db:singleQuery("SELECT id FROM castro_shop_discounts WHERE code = ?", http.postValues.code)

    if check ~= nil then
        session:setFlash("validationError", "Discount code already in use")
        http:redirect("/subtopic/admin/shop/discount/new")
        return
    end

    local validUntil = time:parseDate(http.postValues["until"], "2006-01-02")
    local unlimited = tonumber(http.postValues.unlimited) == 1

    print(unlimited)

    if unlimited then
        http.postValues["code-uses"] = 0
    end

    db:execute(
        "INSERT INTO castro_shop_discounts (code, created_at, valid_till, discount, uses, unlimited) VALUES (?, ?, ?, ?, ?, ?)",
        http.postValues.code,
        os.time(),
        validUntil,
        http.postValues.amount,
        http.postValues["code-uses"],
        http.postValues.unlimited
    )

    session:setFlash("success", "Discount code created")
    http:redirect("/subtopic/admin/shop")
end