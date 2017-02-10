if not session:isLogged() then
    http:redirect("/")
    return
end

local guild = db:singleQuery("SELECT name, ownerid, id FROM guilds WHERE name = ?", http.postValues["guild-name"])

if guild == nil then
    http:redirect("/")
    return
end

local account = session:loggedAccount()
local player = db:singleQuery("SELECT b.id FROM guild_invites a, players b, accounts c WHERE a.guild_id = ? AND a.player_id = b.id AND b.name = ? AND b.account_id = c.id AND c.id = ?", guild.id, http.postValues["character-name"], account.ID)

if player == nil then
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
    return
end


db:execute("DELETE FROM guild_invites WHERE player_id = ? AND guild_id = ?", player.id, guild.id)
session:setFlash("success", "Invitation declined")
http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))