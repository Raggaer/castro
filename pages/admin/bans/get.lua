local config = require("bans")

function get()
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    local data = {}

    local characters = db:query("SELECT id, group_id, name FROM players WHERE account_id = ?", session:loggedAccount().ID)
    if characters then
        local g = xml:unmarshalFile(serverPath .. "/data/XML/groups.xml")
        for _, group in pairs(g.groups.group) do
            if tonumber(group["-access"]) > 0 then
                for _, char in pairs(characters) do
                    if char.group_id == tonumber(group["-id"]) then
                        data.admin = char.id
                        break
                    end
                end
            end
        end
    end

    if not data.admin then
        data.validationError = string.format("You must have a character in an access enabled group in order to issue banishments.")
        http:render("bans.html", data)
        return
    end

    data.success = session:getFlash("success")
    data.validationError = session:getFlash("validationError")

    data.account_bans, ab_cached = db:query("SELECT * FROM account_bans", true)
    data.ip_bans, ib_cached = db:query("SELECT * FROM ip_bans", true)
    data.ban_history, bh_cached = db:query("SELECT * FROM account_ban_history", true)
    data.namelocks, nl_cached = db:query("SELECT * FROM player_namelocks", true)
    data.quickBans = config.quickList

    if data.account_bans and not ab_cached then
        for _, b in pairs(data.account_bans) do
            b.name = db:singleQuery("SELECT name FROM accounts WHERE id = ?", b.account_id).name
            b.start_date = time:parseUnix(tonumber(b.banned_at)).Result
            b.end_date = time:parseUnix(tonumber(b.expires_at)).Result
            b.banned_by_name = db:singleQuery("SELECT name FROM players WHERE id = ?", b.banned_by).name
        end
    end

    if data.ip_bans and not ib_cached then
        for _, b in pairs(data.ip_bans) do
            b.start_date = time:parseUnix(tonumber(b.banned_at)).Result
            b.end_date = time:parseUnix(tonumber(b.expires_at)).Result
            b.banned_by_name = db:singleQuery("SELECT name FROM players WHERE id = ?", b.banned_by).name
        end
    end

    if data.ban_history and not bh_cached then
        for _, b in pairs(data.ban_history) do
            b.name = db:singleQuery("SELECT name FROM accounts WHERE id = ?", b.account_id).name
            b.start_date = time:parseUnix(tonumber(b.banned_at)).Result
            b.end_date = time:parseUnix(tonumber(b.expired_at)).Result
            b.banned_by_name = db:singleQuery("SELECT name FROM players WHERE id = ?", b.banned_by).name
        end
    end

    if data.namelocks and not nl_cached then
        for _, n in pairs(data.namelocks) do
            n.name = db:singleQuery("SELECT name FROM players WHERE id = ?", n.player_id).name
            n.date = time:parseUnix(tonumber(n.namelocked_at)).Result
            n.namelocked_by_name = db:singleQuery("SELECT name FROM players WHERE id = ?", n.namelocked_by).name
        end
    end

    http:render("bans.html", data)
end
