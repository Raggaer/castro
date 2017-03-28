function get()
    if not app.Custom.Forum.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data.list = db:query("SELECT id, title, description FROM castro_forum_category ORDER BY id")

    if data.list ~= nil then
        for i, forum in pairs(data.list) do
            data.list[i].count = db:singleQuery("SELECT COUNT(1) as c FROM castro_forum_post WHERE category_id = ?", forum.id)
            data.list[i].last = db:singleQuery("SELECT a.title as title, b.name as name, a.created_at as created FROM castro_forum_post a, players b WHERE a.author = b.id AND a.category_id = ? ORDER BY a.id DESC", forum.id)
            if data.list[i].last ~= nil then
                if data.list[i].last.title:len() > 12 then
                    data.list[i].last.title = string.sub(data.list[i].last.title, 0, 12) .. "..."
                end
                data.list[i].last.created = time:parseUnix(data.list[i].last["created"])
            end
        end
    end

    http:render("forums.html", data)
end