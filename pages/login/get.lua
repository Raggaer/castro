function get()
    if session:isLogged() then
        http:redirect("/")
        return
    end

    local data = {}

    data["validationError"] = session:getFlash("validationError")
    data["success"] = session:getFlash("success")

    http:render("login.html", data)
end