if session:isLogged() then
    http:redirect("/")
    return
end

local data = {}

data["serverName"] = config:get("ServerName")
data["validationError"] = session:getFlash("validationError")

http:render("register.html", data)

