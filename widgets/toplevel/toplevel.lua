function widget()

    local data = {}

    data.top = db:query("SELECT name, level FROM players ORDER BY level DESC LIMIT 5", true)

    widgets:render("toplevel.html", data)

end

