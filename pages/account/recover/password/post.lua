function post()
    if session:isLogged() then
        http:redirect("/")
        return
    end

    if app.Captcha.Enabled then
        if not captcha:verify(http.postValues["g-recaptcha-response"]) then
            session:setFlash("validationError", "Invalid captcha answer")
            http:redirect("/subtopic/account/recover/password")
            return
        end
    end

    local account = db:singleQuery(
        "SELECT id, name, email FROM accounts WHERE email = ? AND name = ?",
        http.postValues["email"],
        http.postValues["name"]
    )

    if account == nil then
        session:setFlash("validationError", "Wrong email or account name")
        http:redirect("/subtopic/account/recover/password")
        return
    end

    local pp = crypto:randomString(8)

    db:execute("UPDATE accounts SET password = ?  WHERE id = ?", crypto:sha1(pp), account.id)

    if app.Mail.Enabled then
        events:new(
            function()
                local m = {}
                m.to = account.email
                m.subject = "Account password recovery"
                m.body = "Your new account password is: <b>" .. pp .. "</b>"
                mail:send(m)
            end
        )
        session:setFlash("success", "Account found. You will get an email with your password soon")
        http:redirect("/subtopic/account/recover/password")
        return
    end

    session:setFlash("success", "Account found. Your new password is: " .. pp)
    http:redirect("/subtopic/account/recover/password")
end