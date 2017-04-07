require "bbcode"
require "util"

function get()
    -- Block access for anyone who is not admin
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    if not app.Plugin.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data.info = json:unmarshal(http:get(app.Plugin.Origin .. "/plugin/view/" .. http.getValues.id))

    if data.info.Error then
        http:redirect("/subtopic/extensions/list")
        return
    end

    data.info.Description = data.info.Description:parseBBCode()
    data.info.Type = pluginTypeToString(data.info.Type)

    http:render("viewextension.html", data)
end