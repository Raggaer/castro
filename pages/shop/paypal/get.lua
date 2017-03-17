require "paypal"

function get()
    if not app.PayPal.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data["validationError"] = session:getFlash("validationError")
    data["success"] = session:getFlash("success")

    data.validationError = session:getFlash("validationError")
    data.currency = app.PayPal.Currency
    data.list = paypalList
    data.logged = session:isLogged()

    http:render("paypal.html", data)
end