function get()
    local data = {}
    local name = url:decode(http.getValues.name)

    data.info, cache = db:singleQuery("SELECT d.name AS rank, c.name AS guild, a.name, a.stamina, a.health, a.healthmax, a.mana, a.manamax, a.sex, a.vocation, a.level, a.town_id, a.lastlogin, a.lastlogout, a.maglevel, a.skill_sword, a.skill_axe, a.skill_club, a.skill_dist, a.skill_fist, a.skill_shielding, a.skill_fishing FROM players a LEFT JOIN guild_membership b ON b.player_id = a.id LEFT JOIN guilds c ON c.id = b.guild_id LEFT JOIN guild_ranks d ON d.id = b.rank_id WHERE a.name = ?", name, true)

    if data.info == nil then
        http:redirect("/")
        return
    end

    if not cache then
        data.info.vocation = xml:vocationByID(data.info.vocation)
        data.info.town = otbm:townByID(data.info.town_id)
        data.info.lastlogin = time:parseUnix(data.info.lastlogin)
        data.info.lastlogout = time:parseUnix(data.info.lastlogout)
        data.info.healthpercent = (data.info.health * 100) / data.info.healthmax

        if data.info.manamax == 0 then
            data.info.manapercent = 100
        else
            data.info.manapercent = (data.info.mana * 100) / data.info.manamax
        end
    end

    http:render("viewcharacter.html", data)
end