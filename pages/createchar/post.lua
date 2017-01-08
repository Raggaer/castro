if not session:isLogged() then
    http:redirect("/")
    return
end

if not validator:validVocation(http.postValues["character-vocation"]) then
    session:setFlash("validation-error", "Invalid character vocation. Vocation not found")
    http:redirect("/subtopic/createchar")
end

if not validator:validTown(http.postValues["character-town"]) then
    session:setFlash("validation-error", "Invalid character town. Town not found")
    http:redirect("/subtopic/createchar")
end

if not validator:validUsername(http.postValues["character-name"]) then
    session:setFlash("validation-error", "Invalid character name format. Only letters A-Z and spaces allowed")
    http:redirect("/subtopic/createchar")
end

if db:query("SELECT id FROM players WHERE name = ?", http.postValues["character-name"]) ~= nil then
    session:setFlash("validation-error", "Character name already in use")
    http:redirect("/subtopic/createchar")
end

account = session:loggedAccount()

db:execute("INSERT INTO players (name, account_id, vocation, town_id, conditions) VALUES (?, ?, ?, ?, '')", http.postValues["character-name"], account.id, xml:vocationByName(http.postValues["character-vocation"]).id, otbm:townByName(http.postValues["character-town"]).id)