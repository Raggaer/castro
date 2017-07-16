require "paginator"

function get()
    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/library/houses")
        return
    end

    local data = {}

    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect()
        return
    end

    local houseList = {}

    if http.getValues.town == nil then
        houseList = otbm:houseList()
    else
        houseList = otbm:houseList(tonumber(http.getValues.town))
        data.current = otbm:townByID(tonumber(http.getValues.town))
        data.townId = http.getValues.town
    end

    local totalLength = 0

    for _ in pairs(houseList) do
        totalLength = totalLength + 1
    end

    local pg = paginator(page, 15, tonumber(totalLength))

    local newTotalLength = 0

    data.list = {}

    for k in pairs(houseList) do
        if k >= pg.offset and newTotalLength < pg.limit then
            data.list[newTotalLength] = houseList[k]
            newTotalLength = newTotalLength + 1
        end
    end

    data.paginator = pg
    data.towns = otbm:townList()

    http:render("houselist.html", data)
end