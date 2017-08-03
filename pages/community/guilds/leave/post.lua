require "guild"

function post()
    if not session:isLogged() then
        http:redirect("/")
        return
    end

    local guild = db:singleQuery("SELECT id, name, ownerid, id FROM guilds WHERE name = ?", http.postValues["guild-name"])

    if guild == nil then
        http:redirect("/")
        return
    end

    if isGuildOwner(session:loggedAccount().ID, guild) then
        http:redirect("/")
        return
    end

    local player = db:singleQuery("SELECT id, account_id FROM players WHERE name = ?", http.postValues["character-name"])

    if player == nil then
        http:redirect("/")
        return
    end

    if tonumber(player.account_id) ~= session:loggedAccount().ID then
         http:redirect("/")
         return
    end

    if not isGuildMember(player.id, guild) then
        http:redirect("/")
        return
    end

    db:execute("DELETE FROM guild_membership WHERE player_id = ? AND guild_id = ?", player.id, guild.id)

     session:setFlash("success", "You left the guild")
     http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
end