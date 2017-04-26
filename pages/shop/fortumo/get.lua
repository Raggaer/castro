function get()
    if not app.Fortumo.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local data = {}

    data.accountId = session:loggedAccount().Name
    data.serviceId = app.Fortumo.Service
    data.validationError = session:getFlash("validationError")
    data.success = session:getFlash("success")

    http:render("fortumo.html", data)
end