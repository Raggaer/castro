if not session:isLogged() then
    http:redirect("/subtopic/login")
    return
end

session:destroy()
session:setFlash("success", "Logged out")

http:redirect("/subtopic/login")