if not session:isLogged() then
    http:redirect("/")
    return
end

if not validator:validUsername(http.postValues["character-name"]) then
    session:setFlash("validation-error", "Invalid character name format. Only letters A-Z and spaces allowed")
    http:redirect("/subtopic/createchar")
end
