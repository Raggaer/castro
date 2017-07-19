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

    local character = db:singleQuery("SELECT c.id as rankId, c.level, a.id FROM players a, guild_membership b, guild_ranks c WHERE c.id = b.rank_id AND b.player_id = a.id AND b.guild_id = ? AND a.name = ?", guild.id, http.postValues["leader-name"])

    if character == nil then
        http:redirect("/")
        return
    end

    if tonumber(character.level) == 3 then
        http:redirect("/")
        return
    end

    local leaderRank = db:singleQuery("SELECT id FROM guild_ranks WHERE level = 3 AND guild_id = ?", guild.id)
    local memberRank = db:singleQuery("SELECT id FROM guild_ranks WHERE level = 1 AND guild_id = ?", guild.id)

    if memberRank == nil then
        http:redirect("/")
        return
    end

    if leaderRank == nil then
        http:redirect("/")
        return
    end

    db:execute("UPDATE guild_membership SET rank_id = ? WHERE player_id = ? AND guild_id = ?", memberRank.id, guild.ownerid, guild.id)
    db:execute("UPDATE guilds SET ownerid = ? WHERE id = ?", character.id, guild.id)
    db:execute("UPDATE guild_membership SET rank_id = ? WHERE player_id = ? AND guild_id = ?", leaderRank.id, character.id, guild.id)

    session:setFlash("success", "Guild leadership updated")

    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
end