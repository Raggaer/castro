function get()
    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    local account = session:loggedAccount()
    local data = {}

    data.success = session:getFlash("success")
    data.validationError = session:getFlash("validationError")
    data.list = db:query("SELECT name, vocation, level, deletion FROM players WHERE account_id = ? ORDER BY id DESC", account.ID)
    data.account = session:loggedAccount()
    data.account.Creation = time:parseUnix(data.account.Creation)
    data.account.Lastday = time:parseUnix(data.account.Lastday)

    if data.list then
        for _, character in pairs(data.list) do
            character.deletion = tonumber(character.deletion)
            if character.deletion ~= 0 then
                if character.deletion > os.time() then
                    data.account.PendingDeletions = data.account.PendingDeletions or {}
                    table.insert(data.account.PendingDeletions, {name = character.name, deletion = time:parseUnix(character.deletion)})
                end
            end
        end
    end

    if account.Secret == nil then
        data.twofa = false
    else
        data.twofa = true
    end

    if data.account.Premdays > 0 then
        data.account.IsPremium = true
    end

    http:render("dashboard.html", data)
end