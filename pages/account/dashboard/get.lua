function get()
if not session:isLogged() then
    http:redirect("/subtopic/login")
    return
end

local data = {}

data.success = session:getFlash("success")
data.list = db:query("SELECT name, vocation, level FROM players WHERE account_id = ? ORDER BY id DESC", session:loggedAccount().ID)
data.account = session:loggedAccount()
data.account.Creation = time:parseUnix(data.account.Creation)
data.account.Lastday = time:parseUnix(data.account.Lastday)

if data.account.Premdays > 0 then
    data.account.IsPremium = true
end

http:render("dashboard.html", data)
end