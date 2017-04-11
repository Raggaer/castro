-- This file will be executed at start-up

if app.Mode == "dev" then
    print(">> Running on development mode. Never have development mode open to the public")
end

if app.Custom.OnlineChart.Enabled then
	events:add("engine/onlinechart.lua", app.Custom.OnlineChart.Interval)
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
