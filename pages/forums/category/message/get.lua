require "bbcode"
require "paginator"

function get()
    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/index")
        return
    end

    local data = {}
    local messageCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_forum_message WHERE post_id = ?", http.getValues.id)
    local pg = paginator(page, tonumber(app.Custom.Forum.MessagesPerThread), tonumber(messageCount.total))

    data["validationError"] = session:getFlash("validationError")
    data["success"] = session:getFlash("success")
    data.paginator = pg
    data.logged = session:isLogged()
    data.info = db:singleQuery("SELECT id, title FROM castro_forum_post WHERE id = ?", http.getValues.id)

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    data.list = db:query("SELECT a.id as id, a.message as msg, b.level as level, b.vocation as voc, b.name as name, a.created_at as created FROM castro_forum_message a, players b WHERE a.author = b.id AND a.post_id = ? ORDER BY a.created_at LIMIT ?, ?",  http.getValues.id, pg.limit, pg.offset)

    if data.list ~= nil then
        for i, msg in pairs(data.list) do
            data.list[i].created = time:parseUnix(msg.created)
            data.list[i].voc = xml:vocationByID(msg.voc)
            data.list[i].msg = data.list[i].msg:parseBBCode()
        end
    end

    http:render("message.html", data)
end