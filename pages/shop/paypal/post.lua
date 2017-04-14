require "paypal"
require "util"

function post()
    if not app.PayPal.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local package = paypalList[http.postValues.pkg]

    if package == nil then
        http:redirect()
        return
    end

    --[[
        Package name (for paypal description)
        Package price (amount to pay)
        Payment custom field (used later to give points)
        Cancel URL (where to redirect if user cancels the payment process)
        Return URL (where to redirect when user approves the payment)
    ]]--

    local info = paypal:createPayment(
        package.name,
        package.price,
        session:loggedAccount().Name,
        runningURL() .. "/subtopic/shop/paypal",
        runningURL() .. "/subtopic/shop/paypal/review"
    )

    db:execute("INSERT INTO castro_paypal_payments (package_name, state, payment_id, custom, created_at) VALUES (?, ?, ?, ?, ?)", info.Name, info.State, info.PaymentID, info.Custom, os.time())
    http:redirect(info.Link)
end