function widget()
    local data = {}

    data.top = db:query("SELECT name, level, group_id FROM players WHERE group_id < ? ORDER BY level DESC LIMIT 5", app.Custom.HighscoreIgnoreGroup)

    widgets:render("toplevel.html", data)
end

