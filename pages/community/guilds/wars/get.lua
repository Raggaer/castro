require "guild"

function get()
    -- War list
    local data = {}

    data.validationError = session:getFlash("validationError")
    data.list = getActiveGuildWars()

    http:render("guildwarslist.html", data)
end
