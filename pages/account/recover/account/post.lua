function post()
    if session:isLogged() then
        http:redirect("/")
        return
    end

    if app.Captcha.Enabled then
        if not captcha:verify(http.postValues["g-recaptcha-response"]) then
            session:setFlash("validationError", "Invalid captcha answer")
            http:redirect("/subtopic/account/recover/account")
            return
        end
    end

    local account = db:singleQuery(
        "SELECT email, name FROM accounts WHERE email = ? AND password = ?",
        http.postValues["email"],
        crypto:sha1(http.postValues["password"])
    )

    if account == nil then
        session:setFlash("validationError", "Wrong email or password")
        http:redirect("/subtopic/account/recover/account")
        return
    end

    if app.Mail.Enabled then
        events:new(
            function()
                local m = {}
                m.to = account.email
                m.subject = "Account name recovery"
                m.body = "Your account name is: <b>" .. account.name .. "</b>"
                mail:send(m)
            end
        )
        session:setFlash("success", "Account found. You will get an email with your username soon")
        http:redirect("/subtopic/account/recover/account")
        return
    end

    session:setFlash("success", "Account found. Your account name is: " .. account.name)
    http:redirect("/subtopic/account/recover/account")
end