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
        http:redirect("/subtopic/account/fortumo")
        return
    end

    local account = session:loggedAccount()
    local paymentCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_fortumo_payments WHERE account = ?", account.Name)
    local pg = paginator(page, 15, tonumber(paymentCount.total))
    local data = {}

    data.list = db:query("SELECT id, payment_id AS name, created_at AS created, price, points FROM castro_fortumo_payments WHERE account = ? ORDER BY created_at DESC LIMIT ?, ?", account.Name, pg.offset, pg.limit)
    data.paginator = pg

    if data.list ~= nil then
        for i, payment in pairs(data.list) do
            data.list[i].created = time:parseUnix(payment.created)
        end
    end

    http:render("fortumopayments.html", data)
end