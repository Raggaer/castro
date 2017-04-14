require "paypal"

function post()
    if not app.PayPal.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local payment = db:singleQuery("SELECT created_at, payment_id, state, payer_id, package_name, custom FROM castro_paypal_payments WHERE id = ? AND custom = ?", http.postValues["id"], session:loggedAccount().Name)

    if payment.created_at + (60*60*3) < os.time() then
        session:setFlash("validationError", "Payment is 3 days old. Please create a new payment")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    if payment == nil then
        session:setFlash("validationError", "Invalid payment")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    if payment.payer_id == nil then
        session:setFlash("validationError", "Invalid payment")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    if payment.state ~= "created" then
        session:setFlash("validationError", "Invalid payment state. Please approve the payment first")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    local pkg = paypalList[payment["package_name"]]

    if pkg == nil then
        session:setFlash("validationError", "Invalid package")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    if paypal:executePayment(payment["payment_id"], payment["payer_id"]) == false then
        session:setFlash("validationError", "Invalid payment")
        http:redirect("/subtopic/shop/paypal")
        return
    end

    db:execute("UPDATE castro_accounts a, accounts b SET a.points = points + ? WHERE a.account_id = b.id AND b.name = ?", pkg.points, payment.custom)
    db:execute("UPDATE castro_paypal_payments SET state = ? WHERE id = ?", "executed", http.postValues["id"])

    session:setFlash("success", "Package " .. pkg.name .. " purchased. " .. pkg.points .. " points given")

    http:redirect("/subtopic/shop/paypal")
end