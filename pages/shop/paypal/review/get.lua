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

    if info == nil then
        http:redirect("/")
        return
    end

    if info.State ~= "created" then
        session:setFlash("validationError", "Invalid payment state. Payment is not created")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    if info.created_at + (60*60*3) < os.time() then
        session:setFlash("validationError", "Payment is 3 days old. Please create a new payment")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    if info.payer_id == nil then
        db:execute("UPDATE castro_paypal_payments SET payer_id = ? WHERE id = ?", http.getValues["PayerID"], info.id)
    end

    local package = paypalList[info.Name]

    if package == nil then
        http:redirect("/")
        return
    end

    local data = {}

    data.pkg = package
    data.id = info.id

    http:render("review.html", data)
end