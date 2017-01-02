if session.logged then
    http:redirect("/")
    return
end

local data = {}

data["serverName"] = config:get("ServerName")

http:render("register.html", data)

