require "guild"

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

    if db:singleQuery("SELECT 1 FROM guild_invites WHERE player_id = ? AND guild_id = ?", invite.id, guild.id) ~= nil then
        session:setFlash("validationError", "Character already invited")
        http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
        return
    end

    if not isGuildOwner(session:loggedAccount().ID, guild) then
        http:redirect("/")
        return
    end

    db:execute("INSERT INTO guild_invites (player_id, guild_id) VALUES (?, ?)", invite.id, guild.id)
    session:setFlash("success", "Invitation sent")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
end