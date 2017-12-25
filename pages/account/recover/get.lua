function get()
    if session:isLogged() then
        http:redirect("/")
        return
    end
    
    http:render("recover.html")
end