require "guild"

function get()
    -- Viewing single war
    local warid, data = tonumber(http.getValues.id), {}

    data.war = db:singleQuery("SELECT *, (SELECT COUNT(*) AS kills FROM guildwar_kills WHERE warid = a.id AND killerguild = a.guild1) AS guild1_kills, (SELECT COUNT(*) AS kills FROM guildwar_kills WHERE warid = a.id AND killerguild = a.guild2) AS guild2_kills FROM guild_wars a WHERE id = ? ORDER BY started", warid)

    if not data.war then
        session:setFlash("validationError", "Could not find any war with specified id.")
        http:redirect("/subtopic/community/guilds/wars")
    end

    if not cached then
        for k, v in pairs(data.war) do
            -- Convert all integers to actual number
            data.war[k] = tonumber(v) or v
        end
        data.war.leader = (data.war.guild1_kills == data.war.guild2_kills) and "Tied" or (data.war.guild1_kills > data.war.guild2_kills) and data.war.name1 or data.war.name2
        data.war.status_name = warStatusToName(data.war.status)
    end

    data.kills, cached = db:query("SELECT * FROM guildwar_kills WHERE warid = ?", warid)

    if data.kills and not cached then
        for _, kill in pairs(data.kills) do
            kill.killerguild, kill.targetguild = getGuildName(kill.killerguild), getGuildName(kill.targetguild)
        end
    end

    http:render("viewguildwar.html", data)
end
