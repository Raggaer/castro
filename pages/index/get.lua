require "paginator"

local page = 0

if http.getValues.page ~= nil then
    page = math.floor(tonumber(http.getValues.page) + 0.5)
end

if page < 0 then
    http:redirect("/subtopic/index")
    return
end

local articleCount = db:query("SELECT COUNT(*) as total FROM articles", true)
local pg = paginator(page, 4, tonumber(articleCount.total))
local data = {}

data.articles = db:query("SELECT title, text, created_at FROM articles ORDER BY id DESC LIMIT ?, ?", pg.limit, pg.offset, false, true)

if data.articles == nil then
    http:redirect("/subtopic/index")
    return
end

data.paginator = pg

http:render("home.html", data)
