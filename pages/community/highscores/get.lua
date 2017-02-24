require "paginator"

function get()

local data = {}

data.vocList = xml:vocationList()
data.vocType = tonumber(http.getValues.voc)

local cache = false
local query = ""

data.orderType = tonumber(http.getValues.order)

if data.orderType == nil then
    data.orderType = 0
end

local page = 0

if data.vocType == nil then
    data.vocType = 0
end

if http.getValues.page ~= nil then
    page = math.floor(tonumber(http.getValues.page) + 0.5)
end

if not validator:validVocation(data.vocType) then
    http:redirect("/")
    return
end

local playerCount = db:singleQuery("SELECT COUNT(*) as total FROM players WHERE vocation = ?", data.vocType, true)

data.paginator = paginator(page, 15, tonumber(playerCount.total))

if data.orderType == 0 then
    data.order = "Level"
    query = "level"
elseif data.orderType == 1 then
    data.order = "Magic Level"
    query = "maglevel"
elseif data.orderType == 2 then
    data.order = "Balance"
    query = "balance"
elseif data.orderType == 3 then
    data.order = "First Fighting"
    query = "skill_fist"
elseif data.orderType == 4 then
    data.order = "Sword Fighting"
    query = "skill_sword"
elseif data.orderType == 5 then
    data.order = "Axe Fighting"
    query = "skill_axe"
elseif data.orderType == 6 then
    data.order = "Club Fighting"
    query = "skill_club"
elseif data.orderType == 7 then
    data.order = "Distance Fighting"
    query = "skill_dist"
elseif data.orderType == 8 then
    data.order = "Shielding"
    query = "skill_shielding"
elseif data.orderType == 9 then
    data.order = "Fishing"
    query = "skill_fishing"
end

data.list, cache = db:query("SELECT name, vocation, " .. query .. " AS value FROM players WHERE vocation = ? ORDER BY value DESC LIMIT ?, ?", data.vocType, data.paginator.limit, data.paginator.offset, true)

if data.list ~= nil then
    if not cache then
        for _, val in pairs(data.list) do
            val.vocation = xml:vocationByID(val.vocation)
        end
    end
end

data.voc = xml:vocationByID(data.vocType)

http:render("highscores.html", data)

    end