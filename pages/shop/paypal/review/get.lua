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

    local info = db:singleQuery("SELECT id, package_name as Name, state as State, payment_id, payer_id, custom, created_at FROM castro_paypal_payments WHERE payment_id = ? AND custom = ?", http.getValues["paymentId"], session:loggedAccount().Name)
    local identifier = 0

    if info == nil then
        info = paypal:paymentInformation(http.getValues["paymentId"])

        identifier = db:execute("INSERT INTO castro_paypal_payments (package_name, state, payment_id, payer_id, custom, created_at) VALUES (?, ?, ?, ?, ?, ?)", info.Name, info.State, info.PaymentID, info.PayerID, info.Custom, os.time())
    end

    if info == nil then
        http:redirect("/")
        return
    end

    if info.State ~= "created" then
        session:setFlash("validationError", "Invalid payment state. Payment is not created")
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

    if identifier == 0 then
        data.id = info.id
    else
        data.id = identifier
    end

    http:render("review.html", data)
end