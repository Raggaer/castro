require "paginator"

function get()
    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/account/paypal")
        return
    end

    local account = session:loggedAccount()
    local paymentCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_paypal_payments WHERE custom = ?", account.Name)
    local pg = paginator(page, 15, tonumber(paymentCount.total))
    local data = {}

    data.list = db:query("SELECT id, payer_id, payment_id, package_name AS name, state, created_at AS created FROM castro_paypal_payments WHERE custom = ? ORDER BY created_at DESC LIMIT ?, ?", account.Name, pg.offset, pg.limit)
    data.paginator = pg

    if data.list ~= nil then
        for i, payment in pairs(data.list) do

            if payment.created + (60*60*3) > os.time() then
                data.list[i].executed = payment.state == "executed"
            else
                data.list[i].executed = true
            end

            data.list[i].approved = payment.payer_id ~= nil

            data.list[i].created = time:parseUnix(payment.created)
        end
    end

    http:render("paypalpayments.html", data)
end