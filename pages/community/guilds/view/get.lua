require "guild"

function get()
    local data = {}

    data.guild = db:singleQuery("SELECT a.id, a.ownerid, a.name as guildname, a.creationdata, a.motd, b.name, (SELECT COUNT(1) FROM guild_membership WHERE a.id = guild_id) AS members, (SELECT COUNT(1) FROM guild_membership c, players_online d WHERE c.player_id = d.player_id AND c.guild_id = a.id) AS onl, (SELECT MAX(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as top, (SELECT MIN(f.level) FROM guild_membership e, players f WHERE f.id = e.player_id AND e.guild_id = a.id) as low FROM guilds a, players b WHERE a.ownerid = b.id AND a.name = ?", url:decode(http.getValues["name"]))

    if data.guild == nil then
        http:redirect("/subtopic/community/guilds/list")
        return
    end

    if file:exists("public/images/guild-images/" .. data.guild.guildname .. ".png") then
        data.logo = "/images/guild-images/" .. data.guild.guildname .. ".png"
    else
        data.logo = "/images/guild-images/default-guild-logo.png"
    end

    local characters

    if session:isLogged() then
        data.owner = isGuildOwner(session:loggedAccount().ID, data.guild)
        data.ownerdata = db:singleQuery("SELECT a.name FROM players a, guilds b WHERE a.id = b.ownerid AND b.id = ?", data.guild.id)
        characters = db:query("SELECT id FROM players WHERE account_id = ?", session:loggedAccount().ID)
        data.myGuildCharacters = db:query("SELECT a.id, a.name FROM players a, accounts b, guild_membership c WHERE c.player_id = a.id AND a.account_id = b.id AND b.id = ? AND c.guild_id = ?", session:loggedAccount().ID, data.guild.id)
    else
        data.owner = false
    end

    data.guild.created = time:parseUnix(tonumber(data.guild.creationdata))
    data.memberlist = db:query("SELECT a.name, a.level, c.name as rank, c.level as ranklevel FROM guild_membership b, players a, guild_ranks c WHERE c.id = b.rank_id AND b.player_id = a.id AND b.guild_id = ? ORDER BY c.level DESC", data.guild.id)
    data.success = session:getFlash("success")
    data.validationError = session:getFlash("validationError")
    data.wars = getWarsByGuild(data.guild.id)

    if data.owner then
        data.ranks = db:query("SELECT name, level FROM guild_ranks WHERE guild_id = ? ORDER BY level DESC", data.guild.id)
    end

    data.invitations = db:query("SELECT a.id, a.name FROM players a, guild_invites b WHERE b.guild_id = ? AND b.player_id = a.id", data.guild.id)

    if session:isLogged() and characters ~= nil then
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
