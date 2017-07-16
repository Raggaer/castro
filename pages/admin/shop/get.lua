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

    data.success = session:getFlash("success")
    data.list = db:query("SELECT id, name, created_at, updated_at FROM castro_shop_categories ORDER BY created_at DESC")

    if data.list ~= nil then
        for i, category in pairs(data.list) do
            data.list[i].created_at = time:parseUnix(tonumber(category.created_at))
        end
    end

    data.codes = db:query("SELECT id, code, valid_till, unlimited FROM castro_shop_discounts ORDER BY created_at DESC")

    if data.codes ~= nil then
        for i, code in pairs(data.codes) do
            data.codes[i].available = os.time() < tonumber(code.valid_till)
            data.codes[i].valid_till = time:parseUnix(tonumber(code.valid_till))
        end
    end

    http:render("adminshop.html", data)
end