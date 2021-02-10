function get()
    local data = {}
    local name = url:decode(http.getValues.name)

    data.info, cache = db:singleQuery("SELECT a.id, a.account_id, e.premium_ends_at, e.creation, d.name AS rank, c.name AS guild, a.name, a.stamina, a.sex, a.vocation, a.level, a.town_id, a.lastlogin, a.lastlogout, a.maglevel, a.skill_sword, a.skill_axe, a.skill_club, a.skill_dist, a.skill_fist, a.skill_shielding, a.skill_fishing FROM players a LEFT JOIN guild_membership b ON b.player_id = a.id LEFT JOIN guilds c ON c.id = b.guild_id LEFT JOIN guild_ranks d ON d.id = b.rank_id LEFT JOIN accounts e ON e.id = a.account_id WHERE a.name = ?", name)

    if data.info == nil then
        http:redirect("/")
        return
    end

    data.deaths = db:query("SELECT d.level, p.name AS victim, d.time, d.is_player, d.killed_by, d.unjustified FROM player_deaths AS d INNER JOIN players AS p ON d.player_id = p.id WHERE p.id = ? ORDER BY time DESC LIMIT ?", data.info.id, app.Custom.CharacterView.Deaths, true)

    if not cache then
        data.info.accountCreation = time:parseUnix(data.info.creation)
        data.info.accountType = data.info.premium_ends_at > os.time() and "Premium account" or "Free account"
        data.info.vocation = xml:vocationByID(data.info.vocation)
        data.info.town = otbm:townByID(data.info.town_id)
        data.info.lastlogin = time:parseUnix(data.info.lastlogin)
        data.info.lastlogout = time:parseUnix(data.info.lastlogout)
    end

    data.characterList = db:query("SELECT a.id, a.name, (SELECT EXISTS ( SELECT 1 FROM players_online WHERE player_id = a.id) ) AS online FROM players a, accounts b WHERE a.account_id = b.id AND b.id = ? AND a.id <> ?", data.info.account_id, data.info.id)

    http:render("viewcharacter.html", data)
end