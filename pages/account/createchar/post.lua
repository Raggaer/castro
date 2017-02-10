if not session:isLogged() then
    http:redirect("/")
    return
end

if http.postValues["character-name"]:len() < 5 or http.postValues["character-name"]:len() > 12 then
    session:setFlash("validation-error", "Invalid character name. Names can only have 5 to 12 characters")
    http:redirect()
    return
end

if not validator:validVocation(http.postValues["character-vocation"], true) then
    session:setFlash("validation-error", "Invalid character vocation. Vocation not found")
    http:redirect()
    return
end

if not validator:validTown(http.postValues["character-town"]) then
    session:setFlash("validation-error", "Invalid character town. Town not found")
    http:redirect()
    return
end

if not validator:validUsername(http.postValues["character-name"]) then
    session:setFlash("validation-error", "Invalid character name format. Only letters A-Z and spaces allowed")
    http:redirect()
    return
end

if db:query("SELECT id FROM players WHERE name = ?", http.postValues["character-name"]) ~= nil then
    session:setFlash("validation-error", "Character name already in use")
    http:redirect()
    return
end

local account = session:loggedAccount()

if db:singleQuery("SELECT COUNT(*) as total FROM players WHERE account_id = ?", account.ID).total > 5 then
    session:setFlash("validation-error", "You can only have 5 characters")
    http:redirect()
    return
end

db:execute("INSERT INTO players (name, account_id, vocation, town_id, conditions) VALUES (?, ?, ?, ?, '')", http.postValues["character-name"], account.ID, xml:vocationByName(http.postValues["character-vocation"]).ID, otbm:townByName(http.postValues["character-town"]).ID)
session:setFlash("success", "Character created")
http:redirect("/subtopic/account/dashboard")