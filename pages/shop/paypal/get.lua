function get()
    local data = {}

    print(app.PayPal.Enabled)

    data.validationError = session:getFlash("validationError")

    http:render("paypal.html", data)
end