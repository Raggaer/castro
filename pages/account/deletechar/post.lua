function post()
    local name = http.postValues["delete-character"]

    if not name or name == "" then
        session:setFlash("validationError", "An error occured.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    local character = db:singleQuery("SELECT id, name, account_id from players WHERE name = ?", name)
    if not character.id then
        session:setFlash("validationError", "Cannot find character for deletion.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    if tonumber(character.account_id) ~= tonumber(session:loggedAccount().ID) then
        session:setFlash("validationError", "You may only delete characters on your own account.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    local online = db:singleQuery("SELECT player_id FROM players_online WHERE player_id = ?", character.id)
    if online then
        session:setFlash("validationError", "The character must be offline first.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    local guild = db:singleQuery("SELECT guild_id FROM guild_membership WHERE player_id = ?", character.id)
    if guild then
        session:setFlash("validationError", "You must leave or disband your guild first.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    local delay = os.time() + app.Custom.CharacterDeletionDelay
    db:execute("UPDATE players SET deletion = ? WHERE id = ?", delay, character.id)
    session:setFlash("success", "Character " .. character.name .. " has been marked for deletion.")
    http:redirect("/subtopic/account/dashboard")
end