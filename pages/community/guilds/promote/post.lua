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

    local character = db:singleQuery("SELECT c.id as rankId, c.level, a.id FROM players a, guild_membership b, guild_ranks c WHERE c.id = b.rank_id AND b.player_id = a.id AND b.guild_id = ? AND a.name = ?", guild.id, http.postValues["character-name"])

    if character == nil then
        http:redirect("/")
        return
    end

    if tonumber(character.level) ~= 1 then
        http:redirect("/")
        return
    end

    local nextRank = db:singleQuery("SELECT id FROM guild_ranks WHERE level = 2 AND guild_id = ?", guild.id)

    if nextRank == nil then
        http:redirect("/")
        return
    end

    db:execute("UPDATE guild_membership SET rank_id = ? WHERE player_id = ? AND guild_id = ?", nextRank.id, character.id, guild.id)

    session:setFlash("success", "Member promoted")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
end