function get()
    if not app.PayGol.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data["validationError"] = session:getFlash("validationError")
    data["success"] = session:getFlash("success")

    http:render("paygol.html", data)
end