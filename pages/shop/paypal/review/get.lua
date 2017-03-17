require "paypal"

function get()
    if not app.PayPal.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local info = cache:get("paypal_payment_" .. http.getValues["paymentId"])

    if info == nil then
        info = paypal:paymentInformation(http.getValues["paymentId"])

        cache:set("paypal_payment_" .. http.getValues["paymentId"], info, "30m")
    end

    if info == nil then
        http:redirect("/")
        return
    end

    if session:loggedAccount().Name ~= info.Custom then
        http:redirect("/")
        return
    end

    if info.State ~= "created" then
        session:setFlash("validationError", "Invalid payment state")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    local package = paypalList[info.Name]

    if package == nil then
        http:redirect("/")
        return
    end

    local data = {}

    data.pkg = package
    data.paymentId = http.getValues["paymentId"]
    data.payerId = http.getValues["PayerID"]

    http:render("review.html", data)
end