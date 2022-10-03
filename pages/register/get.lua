function get()
    if session:isLogged() then
        http:redirect("/")
        return
    end

    local data = {}

    data["serverName"] = config:get("serverName")
    data["validationError"] = session:getFlash("validationError")
    data["captchaEnabled"] = app.Captcha.Enabled
    data["captchaKey"] = app.Captcha.Public

    http:render("register.html", data)
end