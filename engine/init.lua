-- This file will be executed at start-up

if app.Mode == "dev" then
    print(">> Running on development mode. Never have development mode open to the public")
end

if app.Custom.OnlineChart.Enabled then
	events:add("engine/onlinechart.lua", app.Custom.OnlineChart.Interval)
end

if app.Plugin.Enabled then
    local pluginlist = db:query("SELECT name, plugin_id, updated_at FROM castro_extension_subscribe ORDER BY id DESC")

    if pluginlist ~= nil then
        local idlist = ""

        for _, plugin in pairs(pluginlist) do
            if idlist == "" then
                idlist = plugin.plugin_id
            else
                idlist = idlist .. "," .. plugin.plugin_id
            end
        end

        local c = {}

        c.method = "POST"
        c.url = app.Plugin.Origin .. "/plugin/view/multiple"
        c.timeout = "5s"
        c.authentication = {}
        c.authentication.username = app.Plugin.Username
        c.authentication.password = app.Plugin.Password
        c.data = {}
        c.data.id = idlist
        c.headers = {}
        c.headers["Content-Type"] = "application/x-www-form-urlencoded"

        local resp, _, code = http:curl(c)

        if code == 500 then
            print(">> Cannot get plugin updates: " .. json:unmarshal(resp).Message)
        else
            local resp = json:unmarshal(resp)

            for key, plugin in pairs(resp.object) do
                 local last = time:parseDate(plugin.UpdatedAt, "2006-01-02T15:04:05Z")
                 if last > tonumber(pluginlist[key].updated_at) then
                     print(">> Plugin '" .. plugin.Name .. "' needs to be updated")
                 else
                     print(">> Plugin '" .. plugin.Name .. "' up to date")
                 end
            end
        end
    end
end

if app.CheckUpdates and app.Version ~= "" then
    local commitData = json:unmarshal(http:get("https://api.github.com/repos/Raggaer/castro/commits"))
    local outdated = 0

    for k, v in pairs(commitData.object) do
        if tostring(v.sha) == app.Version then
            break
        else
            outdated = outdated + 1
        end
    end

    print(">> Castro started at " .. time:parseUnix(os.time()).Result .. " - Running " .. outdated .. " commits behind")
else
    print(">> Castro started at " .. time:parseUnix(os.time()).Result)
end
