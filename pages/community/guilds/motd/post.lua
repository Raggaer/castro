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

    if not isGuildOwner(session:loggedAccount().ID, guild) then
        http:redirect("/")
        return
    end

    if http.postValues["guild-motd"]:len() > 50 then
        session:setFlash("validationError", "Motd message must be between 0 - 50 characters")
        http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
        return
    end

    db:execute("UPDATE guilds SET motd = ? WHERE id = ?", http.postValues["guild-motd"], guild.id)
    session:setFlash("success", "Motd updated")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
end