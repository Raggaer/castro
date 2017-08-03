require "guild"

function post()
    if not session:isLogged() then
        http:redirect("/")
        return
    end

    local guild = db:singleQuery("SELECT name, ownerid, id FROM guilds WHERE name = ?", http.postValues["guild-name"])

    if guild == nil then
        http:redirect("/")
        return
    end

    if not isGuildOwner(session:loggedAccount().ID, guild) then
        http:redirect("/")
        return
    end

    http:parseMultiPartForm()

    local logoImage = http:formFile("guild-logo")

    if logoImage == nil then
        session:setFlash("validationError", "Invalid logo image")
        http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
        return
    end

    if not logoImage:isValidPNG() then
        session:setFlash("validationError", "Logo image can only be png")
        http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
        return
    end

    logoImage:saveFileAsPNG("public/images/guild-images/" .. guild.name .. ".png", 64, 64)

    session:setFlash("success", "Logo updated")
    http:redirect("/subtopic/community/guilds/view?name=" .. url:encode(guild.name))
end