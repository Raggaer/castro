if not session:isLogged() then
    http:redirect("/subtopic/login")
    return
end

local data = {}

data["list"] = db:query("SELECT name, vocation, level FROM players WHERE account_id = ? ORDER BY id DESC", session:loggedAccount().id)

http:render("dashboard.html", data)
