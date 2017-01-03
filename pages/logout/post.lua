if not session:isLogged() then
    http:redirect("/subtopic/login")
    return
end

session:set("logged", false)
session:setFlash("success", "You were logged out")

http:redirect("/subtopic/login")