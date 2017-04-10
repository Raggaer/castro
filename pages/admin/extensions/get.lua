require "paginator"

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

    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/index")
        return
    end

    local data = {}

    data.error = session:getFlash("Error")

    if http.getValues.name == nil then
        data.list = json:unmarshal(http:get(app.Plugin.Origin .. "/plugin/list?page=" .. page))
    else
        data.author = http.getValues.name
        data.list = json:unmarshal(http:get(app.Plugin.Origin .. "/plugin/list?author=" .. http.getValues.name .. "&page=" .. page))
    end

    if not data.list.Error then
        data.paginator = paginator(page, tonumber(data.list.PerPage), tonumber(data.list.Total))
    end

    data.subs = db:query("SELECT name, plugin_id AS id, updated_at AS up FROM castro_extension_subscribe ORDER BY id DESC")

    if data.subs ~= nil then
        for i, v in pairs(data.subs) do
            data.subs[i].up = time:parseUnix(tonumber(v.up))
        end
    end

    http:render("extensions.html", data)
end