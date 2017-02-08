if not session:isLogged() then
    http:redirect("/subtopic/login")
    return
end

session:destroy()

http:redirect("/subtopic/login")