local data = {}
local cache = false

data.list = db:query("SELECT a.name, b.name as owner, a.creationdata FROM guilds a, players b WHERE a.ownerid = b.id ORDER BY a.creationdata DESC")


for key, val in pairs(data.list) do
    val.creation = time:parseUnix(tonumber(val.creationdata))
end

http:render("guildlist.html", data)

