require "paginator"

local page = 0

if http.getValues.page ~= nil then
	page = math.floor(tonumber(http.getValues.page) + 0.5)
end

if page < 0 then
	http:redirect("/subtopic/community/deaths")
	return
end

local deathCount = db:singleQuery("SELECT COUNT(*) as total FROM player_deaths", true)
local pg = paginator(page, 10, math.floor(tonumber(deathCount.total), 50)) -- Limit to last 50 deaths

local data = {}
data.deaths = db:query("SELECT d.level, p.name AS victim, d.time, d.is_player, d.killed_by, d.unjustified FROM player_deaths AS d INNER JOIN players AS p ON d.player_id = p.id ORDER BY time DESC LIMIT ?, ?", pg.limit, pg.offset, true)

if data.deaths == nil and page > 0 then
	http:redirect("/subtopic/community/deaths")
	return
end

data.paginator = pg

http:render("latestdeaths.html", data)