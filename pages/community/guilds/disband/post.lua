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

    if not isGuildOwner(session:loggedAccount().ID, guild) then
        http:redirect("/")
        return
    end

    if db:singleQuery("SELECT 1 FROM guild_membership WHERE guild_id = ? AND player_id <> ?", guild.id, guild.ownerid) ~= nil then
        session:setFlash("error", "Your guild still has members")
        http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
        return
    end

    db:execute("DELETE FROM guilds WHERE id = ?", guild.id)

    session:setFlash("success", "Guild disbanded")

    http:redirect("/subtopic/community/guilds/list")
end