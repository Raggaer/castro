function widget()
    local data = {}

    data.top = db:query("SELECT name, level, group_id FROM players WHERE group_id = 1 ORDER BY level DESC LIMIT 5")

    widgets:render("toplevel.html", data)
end

