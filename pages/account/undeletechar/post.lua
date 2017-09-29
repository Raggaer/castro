function post()
    local name = http.postValues["undelete-character"]

    if not name or name == "" then
        session:setFlash("validationError", "An error occured.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    local character = db:singleQuery("SELECT id, name from players WHERE name = ?", name)
    if not character.id then
        session:setFlash("validationError", "Cannot find character to undelete.")
        http:redirect("/subtopic/account/dashboard")
        return
    end

    db:execute("UPDATE players SET deletion = ? WHERE id = ?", 0, character.id)
    session:setFlash("success", "Deletion of the character " .. character.name .. " has been cancelled.")
    http:redirect("/subtopic/account/dashboard")
end