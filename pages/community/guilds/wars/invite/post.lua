require "guild"

function post()
    local guild1 = tonumber(http.postValues.guild1)
    local name1 = getGuildName(guild1)

    if http.postValues.invite then
        -- Send invite
        local targetGuild = db:singleQuery("SELECT id, name FROM guilds WHERE name LIKE ? LIMIT 1", url:decode(http.postValues["invite-guild"]))
        if not targetGuild then
            session:setFlash("validationError", "Could not find any guild with that name.")
            http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(name1))
            return
        end

        if targetGuild.id == guild1 then
            session:setFlash("validationError", "You can not declare war against yourselves, silly.")
            http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(name1))
            return
        end

        local guild2, name2 = tonumber(targetGuild.id), targetGuild.name
        local existingWar = db:singleQuery("SELECT status FROM guild_wars WHERE (guild1 = ? OR guild2 = ?) AND (guild1 = ? OR guild2 = ?) AND status < 2 LIMIT 1", guild1, guild1, guild2, guild2)
        if existingWar then
            local status = tonumber(existingWar.status)
            local s = status == 1 and "an active war" or status == 0 and "a pending invitation"

            session:setFlash("validationError", string.format("There is already %s between your guild and %s.", s, name2))
            http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(name1))
            return
        end

        session:setFlash("success", "Your war invitation has been sent.")
        db:execute("INSERT INTO guild_wars (guild1, guild2, name1, name2, status, started, ended) VALUES (?, ?, ?, ?, ?, ?, ?)", guild1, guild2, name1, name2, 0, os.time(), 0)
    else
        local status, warid = 0, nil
        if http.postValues.accept then
            -- Accept
            status, warid = 1, tonumber(http.postValues.accept)
            session:setFlash("success", "The war invitation has been accepted.")
            -- Force cache update
            cache:delete("active-guildwars")
        elseif http.postValues.reject then
            -- Reject
            status, warid = 2, tonumber(http.postValues.reject)
            session:setFlash("success", "The war invitation has been rejected.")
        elseif http.postValues.cancel then
            -- Cancel
            status, warid = 3, tonumber(http.postValues.cancel)
            session:setFlash("success", "The war invitation has been canceled.")
        end
        -- Update war status
        db:execute("UPDATE guild_wars SET status = ? WHERE id = ?", status, warid)
    end

    -- Force cache update on success
    cache:delete("guildwars-" .. tostring(guild1))

    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(name1))
end
