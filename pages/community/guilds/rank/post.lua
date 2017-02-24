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

local characters = db:query("SELECT id FROM players WHERE account_id = ?", session:loggedAccount().ID)
local owner = false

for _, val in pairs(characters) do
    if val.id == tonumber(guild.ownerid) then
        owner = true
        break
    end
end

if not owner then
    http:redirect("/")
    return
end

if not validator:validGuildRank(http.postValues["rank-3"]) or not validator:validGuildRank(http.postValues["rank-2"]) or not validator:validGuildRank(http.postValues["rank-1"]) then
    session:setFlash("validationError", "Invalid rank title. Titles can only contain (A-Z, -) and must have between 5 and 20 characters")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
    return
end

db:execute("UPDATE guild_ranks SET name = ( CASE WHEN level = 3 THEN ? WHEN level = 2 THEN ? WHEN level = 1 THEN ? END ) WHERE guild_id = ?", http.postValues["rank-3"], http.postValues["rank-2"], http.postValues["rank-1"], guild.id)
session:setFlash("success", "Ranks updated")
http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
    end