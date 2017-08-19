-- This file will be executed at start-up

require "extensionhooks"

if app.Mode == "dev" then
    print(">> Running on development mode. Never have development mode open to the public")
end

if app.Custom.OnlineChart.Enabled then
	events:new(
       function()
           local interval = time:parseDuration(app.Custom.OnlineChart.Interval)
           while true do
               local result = db:singleQuery("SELECT COUNT(*) AS count FROM players_online")
               local count, time = result.count, os.time()
               db:execute("INSERT INTO castro_onlinechart (count, time) VALUES (?, ?)", count, time)

               local old = time - (interval * (app.Custom.OnlineChart.Display + 1))
               db:execute("DELETE FROM castro_onlinechart WHERE time < ?", old)

               sleep(app.Custom.OnlineChart.Interval)
           end
       end
    )
end

-- Run extensions onStartup event
executeHook("onStartup")

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
    log:info("Castro started at " .. time:parseUnix(os.time()).Result .. " - Running " .. outdated .. " commits behind")
else
    print(">> Castro started at " .. time:parseUnix(os.time()).Result)
    log:info("Castro started at " .. time:parseUnix(os.time()).Result)
end
