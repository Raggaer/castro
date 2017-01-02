if session:isLogged() then
    http:redirect("/")
    return
end

if not validator:validate("IsEmail", http.postValues.email) then
    session:setFlash("validationError", "Invalid email format")
    http:redirect("/subtopic/register")
    return
end
