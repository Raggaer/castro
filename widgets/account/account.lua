function widget()

    local data = {}

    data["logged"] = session:isLogged()

    if data["logged"] then
        data["account"] = session:loggedAccount()
    end

    widgets:render("account.html", data)

end