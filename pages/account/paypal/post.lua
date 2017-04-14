require "paypal"

function post()
    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local info = db:singleQuery("SELECT payment_id AS payment, id, package_name AS name FROM castro_paypal_payments WHERE custom = ? AND id = ? AND state = ?", session:loggedAccount().Name, http.postValues.id, "created")

    if info == nil then
        http:redirect("/")
        return
    end

    local package = paypalList[info.name]

    if package == nil then
        http:redirect("/")
        return
    end

    http:redirect("/subtopic/shop/paypal/review?paymentId=" .. info.payment)
end