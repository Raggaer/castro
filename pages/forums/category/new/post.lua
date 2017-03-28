require "bbcode"

function post()
    if not session:isLogged() then
        http:redirect("/subtopic/forums")
        return
    end

    local data = {}
    local account = session:loggedAccount()

    data.logged = session:isLogged()
    data.info = db:singleQuery("SELECT id, title, description FROM castro_forum_category WHERE id = ?", http.getValues.id)
    data.characters = db:query("SELECT name FROM players WHERE account_id = ?", account.ID)

    if data.info == nil then
        http:redirect("/subtopic/forums")
        return
    end

    if http.postValues.action == "preview" then
        data.msg = {}
        data.msg.title = http.postValues.title
        data.msg.text = http.postValues.text
        data.preview = http.postValues.text:parseBBCode()
        http:render("newcategory.html", data)
        return
    end

    if http.postValues.action ~= "send" then
        http:redirect("/subtopic/forums")
        return
    end

    if session:get("forum_cd") ~= nil then
        if session:get("forum_cd") > os.time() then
            session:setFlash("validationError", "Please wait " .. session:get("forum_cd") - os.time() .. " seconds")
            http:redirect("/subtopic/forums/category?forum=" .. url:encode(data.info.title))
            return
        end
    end

    local character = db:singleQuery("SELECT id FROM players WHERE account_id = ? AND name = ?", account.ID, http.postValues.char)

    if character == nil then
        http:redirect("/subtopic/forums")
        return
    end

    if http.postValues.title:len() < 4 then
        session:setFlash("validationError", "Thread title too short (4 characters or more)")
        http:redirect("/subtopic/forums/category/new?id=" .. data.info.id)
        return
    end

    if http.postValues.title:len() > 40 then
        session:setFlash("validationError", "Thread title too long (40 characters or less)")
        http:redirect("/subtopic/forums/category/new?id=" .. data.info.id)
        return
    end

    if http.postValues.text:len() < 4 then
        session:setFlash("validationError", "Thread message too short (4 characters or more)")
        http:redirect("/subtopic/forums/category/new?id=" .. data.info.id)
        return
    end

    if http.postValues.title:len() > 250 then
        session:setFlash("validationError", "Thread message too long (250 characters or less)")
        http:redirect("/subtopic/forums/category/new?id=" .. data.info.id)
        return
    end

    local threadid = db:execute(
        "INSERT INTO castro_forum_post (title, category_id, created_at, updated_at, author) VALUES (?, ?, ?, ?, ?)",
        http.postValues.title,
        data.info.id,
        os.time(),
        os.time(),
        character.id
    )

    db:execute(
        "INSERT INTO castro_forum_message (post_id, author, message, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        threadid,
        character.id,
        http.postValues.text,
        os.time(),
        os.time()
    )

    session:set("forum_cd", os.time() + app.Custom.Forum.SpamCooldown)
    session:setFlash("success", "Thread created")
    http:redirect("/subtopic/forums/category/message?id=" .. threadid)
end