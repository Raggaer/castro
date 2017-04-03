function get()
    if not app.PayGol.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local data = {}

    data["validationError"] = session:getFlash("validationError")
    data["success"] = session:getFlash("success")
    data.custom = session:loggedAccount().Name
    data.service = app.PayGol.Service

    http:render("paygol.html", data)
end