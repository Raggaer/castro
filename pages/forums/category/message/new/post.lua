require "bbcode"
require "paginator"

function post()
    if not app.Custom.Forum.Enabled then
        http:redirect("/")
        return
    end

    if not session:isLogged() then
        http:redirect("/subtopic/forums")
        return
    end

    local account = session:loggedAccount()
    local data = {}

    data.info = db:singleQuery("SELECT id, title FROM castro_forum_post WHERE id = ?", http.getValues.id)

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    if session:get("forum_cd") ~= nil then
        if session:get("forum_cd") > os.time() then
            session:setFlash("validationError", "Please wait " .. session:get("forum_cd") - os.time() .. " seconds")
            http:redirect("/subtopic/forums/category/message?id=" .. data.info.id)
            return
        end
    end

    if http.postValues.action == "preview" then
        data.characters = db:query("SELECT name, vocation, level FROM players WHERE account_id = ? ORDER BY id DESC", account.ID)
        data.msg = http.postValues.text
        data.preview = http.postValues.text:parseBBCode()
        http:render("newmessage.html", data)
        return
    end

    if http.postValues.action ~= "send" then
        http:redirect("/subtopic/forums")
        return
    end

    local character = db:singleQuery("SELECT id FROM players WHERE name = ? AND account_id = ?", url:decode(http.postValues.char), account.ID)

    if character == nil then
        http:redirect("/subtopic/forums")
        return
    end

    if http.postValues.text:len() <= 4 then
        session:setFlash("validationError", "Post message too short (4 characters or more)")
        http:redirect("/subtopic/forums/category/message/new?id=" .. data.info.id)
        return
    end

    if http.postValues.text:len() > 255 then
        session:setFlash("validationError", "Post message too large (255 characters or less)")
        http:redirect("/subtopic/forums/category/message/new?id=" .. data.info.id)
        return
    end

    local msgid = db:execute(
        "INSERT INTO castro_forum_message (post_id, message, author, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        http.getValues.id,
        http.postValues.text,
        character.id,
        os.time(),
        os.time()
    )

    local messageCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_forum_message WHERE post_id = ?", http.getValues.id)

    session:set("forum_cd", os.time() + app.Custom.Forum.SpamCooldown)
    session:setFlash("success", "Message created")
    http:redirect("/subtopic/forums/category/message?id=" .. data.info.id .. "&page=" ..  math.floor(messageCount.total / app.Custom.Forum.MessagesPerThread) .. "#message-" .. msgid)
end