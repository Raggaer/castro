local config = require("bans")

function post()
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    local ban = {}
    local player = db:singleQuery("SELECT id, name, account_id, lastip FROM players WHERE name = ?", http.postValues.player_name)
    if not player then
        session:setFlash("validationError", string.format("Could not find player named %q. Note that while not case sensitive the name must be spelled exactly right so we do not ban the wrong person by accident.", http.postValues.player_name))
        http:redirect("/subtopic/admin/bans")
        return
    end

    ban.account_id = player.account_id
    ban.player_id = player.id
    ban.banned_by = tonumber(http.postValues.banned_by)

    if http.postValues.quick_ban then
        local i = tonumber(http.postValues.quick_ban)
        for _, b in pairs(config.quickList) do
            if b.id == i then
                ban.type = b.type
                ban.reason = b.reason
                ban.duration = b.duration and (tonumber(b.duration) or (b.duration ~= "" and time:parseDuration(b.duration)))
                break
            end
        end
    else
        ban.type = http.postValues.ban_type
        -- Accept seconds or duration string
        ban.duration = http.postValues.ban_duration ~= "" and tonumber(http.postValues.ban_duration) or time:parseDuration(http.postValues.ban_duration)
        if not ban.duration then
            session:setFlash("validationError", "Ban duration missing or invalid format.")
            http:redirect("/subtopic/admin/bans")
            return
        end
        ban.reason = http.postValues.ban_reason ~= "" and http.postValues.ban_reason
        if not ban.reason then
            session:setFlash("validationError", "Ban reason can not be empty.")
            http:redirect("/subtopic/admin/bans")
            return
        end
    end

    if ban.type == "account_ban" then
        -- Ban account if not already banned
        if not db:singleQuery("SELECT 1 FROM account_bans WHERE account_id = ?", ban.account_id, false) then
            db:execute("INSERT INTO account_bans (account_id, reason, banned_at, expires_at, banned_by) VALUES (?, ?, ?, ?, ?)", ban.account_id, ban.reason, os.time(), os.time() + ban.duration, ban.banned_by)
            session:setFlash("success", string.format("The account of %s have been banned for %s hours.", player.name, math.floor(ban.duration / 60 / 60)))
            -- Clear cache
            cache:delete("SELECT * FROM account_bans")
        else
            session:setFlash("validationError", string.format("The account of %s is already banned.", player.name))
        end
    elseif ban.type == "ip_ban" then
        -- Ban IP if not already banned
        if not db:singleQuery("SELECT 1 FROM ip_bans WHERE account_id = ?", player.lastip, false) then
            db:execute("INSERT INTO ip_bans (ip, reason, banned_at, expires_at, banned_by) VALUES (?, ?, ?, ?, ?)", player.lastip, ban.reason, os.time(), os.time() + ban.duration, ban.banned_by)
            session:setFlash("success", string.format("The ip address of %s have been banned for %s hours.", player.name, math.floor(ban.duration / 60 / 60)))
            -- Clear cache
            cache:delete("SELECT * FROM ip_bans")
        else
            session:setFlash("validationError", string.format("The last known IP of %s is already banned.", player.name))
        end
    elseif ban.type == "namelock" then
        -- Namelock player if not already namelocked
        if not db:singleQuery("SELECT 1 FROM player_namelocks WHERE player_id = ?", ban.player_id, false) then
            db:execute("INSERT INTO player_namelocks (player_id, reason, namelocked_at, namelocked_by) VALUES (?, ?, ?, ?)", ban.player_id, ban.reason, os.time(), ban.banned_by)
            session:setFlash("success", string.format("%s have been namelocked.", player.name))
            -- Clear cache
            cache:delete("SELECT * FROM player_namelocks")
        else
            session:setFlash("validationError", string.format("%s is already name locked.", player.name))
        end
    elseif ban.type == "account_ip_ban" then
        -- Ban both account and IP
        local success = false
        -- Account ban
        if not db:singleQuery("SELECT 1 FROM account_bans WHERE account_id = ?", ban.account_id, false) then
            db:execute("INSERT INTO account_bans (account_id, reason, banned_at, expires_at, banned_by) VALUES (?, ?, ?, ?, ?)", ban.account_id, ban.reason, os.time(), os.time() + ban.duration, ban.banned_by)
            success = true
            -- Clear cache
            cache:delete("SELECT * FROM account_bans")
        end

        -- IP ban
        if not db:singleQuery("SELECT 1 FROM ip_bans WHERE ip = ?", player.lastip, false) then
            db:execute("INSERT INTO ip_bans (ip, reason, banned_at, expires_at, banned_by) VALUES (?, ?, ?, ?, ?)", player.lastip, ban.reason, os.time(), os.time() + ban.duration, ban.banned_by)
            success = true
            -- Clear cache
            cache:delete("SELECT * FROM ip_bans")
        end

        -- Display success if either one of the ban types succeeded
        if success then
            session:setFlash("success", string.format("The account and ip address of %s have been banned for %s hours.", player.name, math.floor(ban.duration / 60 / 60)))
        else
            session:setFlash("validationError", string.format("The account and IP address of %s is already banned.", player.name))
        end
    else
        session:setFlash("validationError", "Unknown 'type' in banishment data.")
    end

    http:redirect("/subtopic/admin/bans")
end
