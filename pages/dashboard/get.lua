if not session:isLogged() then
    http:redirect("/subtopic/login")
    return
end

http:render("dashboard.html")
