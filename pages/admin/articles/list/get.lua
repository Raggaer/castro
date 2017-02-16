-- Block access for anyone who is not admin
if not session:isLogged() or not session:isAdmin() then
    http:redirect("/")
    return
end

require "paginator"

local page = 0

if http.getValues.page ~= nil then
    page = math.floor(tonumber(http.getValues.page) + 0.5)
end

if page < 0 then
    http:redirect("/subtopic/admin/articles/list")
    return
end

local data = {}
data.success = session:getFlash("success")
data.validationError = session:getFlash("validationError")

-- Fetch articles from directly from database
-- User is admin so should be fine and would be strange not to see changes right away
local articleCount = db:singleQuery("SELECT COUNT(*) as total FROM castro_articles", false)
local pg = paginator(page, 10, tonumber(articleCount.total))
data.articles = db:query("SELECT id, title, text, created_at, updated_at FROM castro_articles ORDER BY id DESC LIMIT ?, ?", pg.limit, pg.offset, false)
--data.articles = db:query("SELECT id, title, text, created_at, updated_at FROM castro_articles ORDER BY id DESC LIMIT 1", false)
data.paginator = pg

if data.articles ~= nil then
	for _, article in pairs(data.articles) do
		--string.parseBBCode(article.text)
		article.created_at = time:parseUnix(tonumber(article.created_at))
		if article.updated_at then
			article.updated_at = time:parseUnix(tonumber(article.updated_at))
		else
			article.updated_at = {Result = "Never"}
		end
		--article.created_at = tonumber(article.created_at)
	end
end

http:render("articles.html", data)