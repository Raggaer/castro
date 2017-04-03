function get()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    http:render("extensions.html", nil)
end