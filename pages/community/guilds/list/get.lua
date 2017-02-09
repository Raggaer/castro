local data = {}

data.logged = session:isLogged()

if data.logged then
    data.characters = db:query("SELECT name FROM players WHERE account_id = ?", session:loggedAccount().ID)
end

data.list = db:query("SELECT a.name, b.name as owner, a.creationdata FROM guilds a, players b WHERE a.ownerid = b.id ORDER BY a.creationdata DESC")

if data.list ~= nil then
    for _, val in pairs(data.list) do
        val.creation = time:parseUnix(tonumber(val.creationdata))
    end
end

data["validationError"] = session:getFlash("validationError")

http:render("guildlist.html", data)

