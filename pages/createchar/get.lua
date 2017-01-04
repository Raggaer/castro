if not session:isLogged() then
    http:redirect("/")
    return
end

http:render("createcharacter.html")