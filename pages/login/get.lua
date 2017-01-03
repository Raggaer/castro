if session:isLogged() then
    http:redirect("/")
    return
end

local data = {}

data["validationError"] = session:getFlash("validationError")

http:render("login.html", data)