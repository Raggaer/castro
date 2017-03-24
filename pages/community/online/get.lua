function get()
	local data = {}

	data.list, cached = db:query("SELECT po.player_id as id, p.name, p.level AS level, p.vocation AS vocation FROM players_online AS po INNER JOIN players AS p ON p.id = po.player_id ORDER BY p.name", true)

	data.count = data.list and #data.list or 0

	if not cached then
		for i = 1, data.count do
			data.list[i].vocation = xml:vocationByID(tonumber(data.list[i].vocation))
		end
	end

	if app.Custom.OnlineChart.Enabled == true then
		local chart = cache:get("onlineChart")
		if not chart then
			local template = {
				labels = {},
				datasets = {
					{
			            label = 'Players online',
			            data = {},
			            backgroundColor = 'rgba(0, 140, 186, 0.2)',
			            borderColor = 'rgba(0, 140, 186, 1)',
			            borderWidth = 1,
			        }
			    }
			}

			local records = db:query("SELECT * FROM (SELECT * FROM castro_onlinechart ORDER BY time DESC LIMIT ?) X ORDER BY time ASC", app.Custom.OnlineChart.Display, false)

			if records then
				for index, row in pairs(records) do
					template.labels[index] = os.date("%H:%M", tonumber(row.time))
					template.datasets[1].data[index] = row.count
				end
			end

			chart = json:marshal(template)
			cache:set("onlineChart", chart, app.Custom.OnlineChart.CacheTime)
		end

		data.chart = chart
	end

	http:render("online.html", data)
end
