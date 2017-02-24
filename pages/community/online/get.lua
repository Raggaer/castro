function get()
local data = {}

data.list, cache = db:query("SELECT po.player_id as id, p.name, p.level AS level, p.vocation AS vocation FROM players_online AS po INNER JOIN players AS p ON p.id = po.player_id ORDER BY p.name", true)

data.count = data.list and #data.list or 0

if not cache then
	for i = 1, data.count do
		data.list[i].vocation = xml:vocationByID(tonumber(data.list[i].vocation))
	end
end

http:render("online.html", data)
end