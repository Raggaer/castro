-- Players online chart background task

local interval = time:parseDuration(app.Custom.OnlineChart.Interval)

function run()
	local result = db:singleQuery("SELECT COUNT(*) AS count FROM players_online")
	local count, time = result.count, os.time()
	db:execute("INSERT INTO castro_onlinechart (count, time) VALUES (?, ?)", count, time)

	local old = time - (interval * (app.Custom.OnlineChart.Display + 1))
	db:execute("DELETE FROM castro_onlinechart WHERE time < ?", old)
end
