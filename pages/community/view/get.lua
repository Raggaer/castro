local data = {}
local name = url:decode(http.getValues.name)

data.info, cache = db:query("SELECT name, stamina, health, healthmax, mana, manamax, sex, vocation, level, town_id, lastlogin, lastlogout FROM players WHERE name = ?", name, true, true)

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