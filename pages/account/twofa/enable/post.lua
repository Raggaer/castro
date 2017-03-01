function post()
    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local account = session:loggedAccount()

    if account.Secret ~= nil then
        http:redirect("/")
        return
    end

    local secret = session:getFlash("twofa-key")

    if secret == nil then
        http:redirect("/")
        return
    end

    if not validator:validQRToken(http.postValues.token, secret) then
        session:setFlash("validationError", "Invalid token. Please try again")
        http:redirect()
        return
    end

    session:destroy()
    db:execute("UPDATE accounts SET secret = ? WHERE id = ?", secret, account.ID)
    session:setFlash("success", "Two-factor authentication enabled. Please log-in")
    http:redirect("/subtopic/login")
end
