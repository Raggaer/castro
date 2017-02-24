function post()
if not session:isLogged() then
    http:redirect("/")
    return
end

local guild = db:singleQuery("SELECT name, ownerid, id FROM guilds WHERE name = ?", http.postValues["guild-name"])

if guild == nil then
    http:redirect("/")
    return
end

local invite = db:singleQuery("SELECT id FROM players WHERE name = ?", http.postValues["invite-name"])

if invite == nil then
    session:setFlash("validationError", "Character not found")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
    return
end

if db:singleQuery("SELECT guild_id FROM guild_membership WHERE player_id = ?", invite.id) ~= nil then
    session:setFlash("validationError", "Character already on a guild")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
    return
end

local characters = db:query("SELECT id FROM players WHERE account_id = ?", session:loggedAccount().ID)
local owner = false

if characters ~= nil then
    for _, val in pairs(characters) do
        if val.id == tonumber(guild.ownerid) then
            owner = true
            break
        end
    end
end

if not owner then
    http:redirect("/")
    return
end

db:execute("INSERT INTO guild_invites (player_id, guild_id) VALUES (?, ?)", invite.id, guild.id)
session:setFlash("success", "Invitation sent")
http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
    end