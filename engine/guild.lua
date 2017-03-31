function isGuildOwner(accountid, guild)
    local characters = db:query("SELECT id FROM players WHERE account_id = ?", accountid)

    for _, val in pairs(characters) do
        if val.id == tonumber(guild.ownerid) then
            return true
        end
    end

    return false
end

function getGuildName(guildId)
    local guild = db:singleQuery("SELECT name FROM guilds where id = ?", guildId)
    return guild.name
end

function warStatusToName(status)
    return (status == 0 and "Pending") or (status == 1 and "Active") or (status == 2 and "Rejected") or (status == 3 and "Canceled") or (status == 4 and "Ended") or "Unknown"
end

-- Get all active wars
-- Returns table of active wars from guild_wars table
-- Extra fields per war: guild1_kills, guild2_kills, leader = name of currently leading guild or "Tied"
function getActiveGuildWars()
    local list = cache:get("active-guildwars")

    if list == nil then
        list = false -- We explicitly cache 'false' to avoid repeated queries when no wars are found
        local wars = db:query("SELECT *, (SELECT COUNT(*) AS kills FROM guildwar_kills WHERE warid = a.id AND killerguild = a.guild1) AS guild1_kills, (SELECT COUNT(*) AS kills FROM guildwar_kills WHERE warid = a.id AND killerguild = a.guild2) AS guild2_kills FROM guild_wars a WHERE status = 1 ORDER BY started")
        if wars ~= nil then
            list = {}
            for _, war in pairs(wars) do
                for k, v in pairs(war) do
                    -- Convert all integers to actual number
                    war[k] = tonumber(v) or v
                end
                war.leader = (war.guild1_kills == war.guild2_kills) and "Tied" or (war.guild1_kills > war.guild2_kills) and war.name1 or war.name2
                war.status_name = warStatusToName(war.status)
                list[tostring(war.id)] = war
            end
        end
        cache:set("active-guildwars", list, "15m")
    end

    return list
end

-- Get all wars for a specific guild
-- Return is sorted by status
-- Ex. {["Pending"] = {war1, war2, ...}}
function getWarsByGuild(guildId)
	local list = cache:get("guildwars-" .. tostring(guildId))

    if not list then
        list = db:query("SELECT *, (SELECT COUNT(*) AS kills FROM guildwar_kills WHERE warid = a.id AND killerguild = a.guild1) AS guild1_kills, (SELECT COUNT(*) AS kills FROM guildwar_kills WHERE warid = a.id AND killerguild = a.guild2) AS guild2_kills FROM guild_wars a WHERE (guild1 = ? OR guild2 = ?) ORDER BY started", guildId, guildId, false) or {}
        local wars = {}

        if list then
            for _, war in pairs(list) do
                for k, v in pairs(war) do
                    war[k] = tonumber(v) or v
                end
                war.leader = (war.guild1_kills == war.guild2_kills) and "Tied" or (war.guild1_kills > war.guild2_kills) and war.name1 or war.name2
                war.status_name = warStatusToName(war.status)
                wars[war.status_name] = wars[war.status_name] or {}
                table.insert(wars[war.status_name], war)
            end
        end

        list = wars
        cache:set("guildwars-" .. tostring(guildId), list, "15m")
    end

    return list
end
