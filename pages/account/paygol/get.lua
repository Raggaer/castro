require "paginator"

function get()
    if not app.PayGol.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/account/paygol")
        return
    end

    local account = session:loggedAccount()
    local paymentCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_paygol_payments WHERE custom = ?", account.Name)
    local pg = paginator(page, 15, tonumber(paymentCount.total))
    local data = {}

    data.list = db:query("SELECT id, transaction_id AS name, created_at AS created, price, points FROM castro_paygol_payments WHERE custom = ? ORDER BY created_at DESC LIMIT ?, ?", account.Name, pg.offset, pg.limit)
    data.paginator = pg

    if data.list ~= nil then
        for i, payment in pairs(data.list) do
            data.list[i].created = time:parseUnix(payment.created)
        end
    end

    http:render("paygolpayments.html", data)
end