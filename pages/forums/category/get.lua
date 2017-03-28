function get()
    local data = {}

    data.logged = session:isLogged()
    data.info = db:singleQuery("SELECT title, description FROM castro_forum_category WHERE title = ?", url:decode(http.getValues.forum))

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    data.list = db:query("SELECT a.id as id, a.title as title, b.name as name, a.created_at as created FROM castro_forum_post a, players b, castro_forum_category c WHERE c.title = ? AND c.id = a.category_id AND a.author = b.id", url:decode(http.getValues.forum))

    if data.list ~= nil then
       for i, post in pairs(data.list) do
          data.list[i].created = time:parseUnix(post.created)
       end
    end

    http:render("category.html", data)
end