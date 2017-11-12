function widget()
    local data = {}

    data["logged"] = session:isLogged()

    if data["logged"] then
        data["account"] = session:loggedAccount()
    end

    data.paypal = app.PayPal.Enabled
    data.paygol = app.PayGol.Enabled
    data.fortumo = app.Fortumo.Enabled
    data.shop = app.Shop.Enabled

    widgets:render("account.html", data)
end