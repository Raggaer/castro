require "paginator"

function get()
    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/community/guilds/list")
        return
    end

    local data = {}

    data.logged = session:isLogged()

    if data.logged then
        data.characters = db:query("SELECT name FROM players WHERE account_id = ?", session:loggedAccount().ID)
    end

    local guildTotal = db:singleQuery("SELECT COUNT(*) AS total FROM guilds")

    local pg = paginator(page, 15, tonumber(guildTotal.total))

    data.list = db:query("SELECT a.name, b.name as owner, a.creationdata FROM guilds a, players b WHERE a.ownerid = b.id ORDER BY a.creationdata DESC LIMIT ?, ?", pg.offset, pg.limit)

    if data.list ~= nil then
        for _, val in pairs(data.list) do
            val.creation = time:parseUnix(tonumber(val.creationdata))
        end
    end

    data.paginator = pg
    data["validationError"] = session:getFlash("validationError")
    data.error = session:getFlash("error")
    data.success = session:getFlash("success")

    http:render("guildlist.html", data)
end