require "paypal"

function get()
    if not app.PayPal.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data.validationError = session:getFlash("validationError")
    data.currency = app.PayPal.Currency
    data.list = paypalList

    http:render("paypal.html", data)
end