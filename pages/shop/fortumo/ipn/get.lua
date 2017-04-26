require "util"

function get()
    if not app.Fortumo.Enabled then
        return
    end

    if not app.Mode == "dev" then
        if http:getRemoteAddress() ~= "1.2.3.4" or http:getRemoteAddress() ~= "2.3.4.5" then
            return
        end
    end

    if http.getValues.service_id ~= app.Fortumo.Service then
        return
    end

    local secret = {}

    for method, content in orderedPairs(http.getValues) do
        if method ~= "sig" then
            table.insert(secret, method .. "=" .. content)
        end
    end

    local secretKey = table.concat(secret) .. app.Fortumo.Secret
    local secretKeyHash = crypto:md5(secretKey)

    if secretKeyHash ~= http.getValues.sig then
        return
    end

    db:execute("UPDATE castro_accounts a, accounts b SET a.points = points + ? WHERE a.account_id = b.id AND b.name = ?", http.getValues.amount, http.getValues.cuid)

    db:execute(
        "INSERT INTO castro_fortumo_payments (account, points, price, currency, sender, operator, payment_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        http.getValues.cuid,
        http.getValues.amount,
        http.getValues.price,
        http.getValues.currency,
        http.getValues.sender,
        http.getValues.operator,
        http.getValues.payment_id,
        os.time()
    )
end
