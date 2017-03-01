function get()
    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local account = session:loggedAccount()

    if account.Secret ~= nil then
        http:redirect("/")
        return
    end

    local secretKey = crypto:qrKey()

    print(secretKey)

    local secret = string.format("otpauth://totp/%s:%s?secret=%s&issuer=%s", config:get("serverName"), account.Name, secretKey, config:get("serverName"))
    local data = {}

    data.validationError = session:getFlash("validationError")
    data.url = "data:image/png;base64," .. crypto:qr(secret)

    session:setFlash("twofa-key", secretKey)

    http:render("enabletwofa.html", data)
end