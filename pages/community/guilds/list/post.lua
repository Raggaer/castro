function post()
if not session:isLogged() then
    http:redirect("/")
    return
end

local account = session:loggedAccount()
local character = db:singleQuery("SELECT id, name FROM players WHERE account_id = ? AND name = ?", account.ID, url:decode(http.postValues["guild-owner"]))

if character == nil then
    http:redirect("/")
    http:redirect()
    return
end

if not validator:validGuildName(http.postValues["guild-name"]) then
    session:setFlash("validationError", "Invalid guild name")
    http:redirect()
    return
end

if db:singleQuery("SELECT 1 FROM guild_membership WHERE player_id = ?", character.id) ~= nil then
    session:setFlash("validationError", "Character is member of a guild")
    http:redirect()
    return
end

local guild_id = db:execute("INSERT INTO guilds (name, ownerid, creationdata, motd) VALUES (?, ?, ?, ?)", http.postValues["guild-name"], character.id, os.time(), "Guild leader must edit this text")
local leader_id = db:singleQuery("SELECT id FROM guild_ranks WHERE guild_id = ? AND level = 3", guild_id)

db:execute("INSERT INTO guild_membership (player_id, guild_id, rank_id) VALUES (?, ?, ?)", character.id, guild_id, leader_id.id)

session:setFlash("success", "Guild created")
http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(http.postValues["guild-name"]))
    end