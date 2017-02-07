if not session:isAdmin() then
    http:redirect("/")
end

http:render("logs.html", nil)