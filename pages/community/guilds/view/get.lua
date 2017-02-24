function get()
    local data = {}
    local cache = false

    data.guild, cache = db:singleQuery("SELECT a.id, a.ownerid, a.name as guildname, a.creationdata, a.motd, b.name, (SELECT COUNT(1) FROM guild_membership WHERE a.id = guild_id) AS members, (SELECT COUNT(1) FROM guild_membership c, players_online d WHERE c.player_id = d.player_id AND c.guild_id = a.id) AS onl, (SELECT MAX(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as top, (SELECT MIN(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as low FROM guilds a, players b WHERE a.ownerid = b.id AND a.name = ?", url:decode(http.getValues["name"]), true)

    if data.guild == nil then
        http:redirect("/subtopic/community/guilds/list")
        return
    end

    local characters

    if session:isLogged() then
        characters = db:query("SELECT id FROM players WHERE account_id = ?", session:loggedAccount().ID)
    end

    data.owner = false

    if characters ~= nil then
        for _, val in pairs(characters) do
            if val.id == tonumber(data.guild.ownerid) then
                data.owner = true
                break
            end
        end
    end

    if not cache then
        data.guild.created = time:parseUnix(tonumber(data.guild.creationdata))
    end

    data.memberlist = db:query("SELECT a.name, a.level, c.name as rank FROM guild_membership b, players a, guild_ranks c WHERE c.id = b.rank_id AND b.player_id = a.id AND b.guild_id = ? ORDER BY c.level DESC", data.guild.id)
    data.success = session:getFlash("success")
    data.validationError = session:getFlash("validationError")

    if data.owner then
        data.ranks = db:query("SELECT name, level FROM guild_ranks WHERE guild_id = ? ORDER BY level DESC", data.guild.id)
    end

    data.invitations = db:query("SELECT a.id, a.name FROM players a, guild_invites b WHERE b.guild_id = ? AND b.player_id = a.id", data.guild.id)

    if session:isLogged() then
        if data.invitations ~= nil then
            for _, v in pairs(data.invitations) do
                for _, p in pairs(characters) do
                    if p.id == v.id then
                        v.m = true
                    end
                end
            end
        end
    end

    http:render("viewguild.html", data)
end