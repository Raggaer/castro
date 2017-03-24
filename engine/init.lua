-- This file will be executed at start-up
if app.Mode == "dev" then
    print("Running on development mode. Never have development mode open to the public")
end

if app.Custom.OnlineChart.Enabled then
	events:add("engine/onlinechart.lua", app.Custom.OnlineChart.Interval)
end

print("Castro started at " .. time:parseUnix(os.time()).Result)
