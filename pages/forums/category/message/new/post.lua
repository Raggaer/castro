require "bbcode"

function post()
    if not session:isLogged() then
        http:redirect("/subtopic/forums")
        return
    end

    local data = {}

    data.info = db:singleQuery("SELECT id, title FROM castro_forum_post WHERE id = ?", http.getValues.id)

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    if http.postValues.action == "preview" then
        data.msg = http.postValues.text
        data.preview = http.postValues.text:parseBBCode()
        http:render("newmessage.html", data)
        return
    end

    if http.postValues.action ~= "send" then
        http:redirect("/subtopic/forums")
        return
    end

    local character = db:singleQuery("SELECT id FROM players WHERE name = ?", url:decode(http.postValues.char))

    if character == nil then
        http:redirect("/subtopic/forums")
        return
    end

    db:execute(
        "INSERT INTO castro_forum_message (post_id, message, author, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        http.getValues.id,
        http.postValues.text,
        character.id,
        os.time(),
        os.time()
    )

    http:redirect("/subtopic/forums/category/message?id=" .. data.info.id)
end