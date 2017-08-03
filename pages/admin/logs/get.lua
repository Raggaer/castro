require "util"
require "paginator"

function get()
    if not session:isAdmin() then
        http:redirect("/")
        return
    end

    local page = 0

    if http.getValues.page ~= nil then
        page = math.floor(tonumber(http.getValues.page) + 0.5)
    end

    if page < 0 then
        http:redirect("/subtopic/admin/logs")
        return
    end

    local data = {}
    local lines = 0
    local file = io.open("logs/" .. logFile, "r")
    local linesData = {}

    for line in file:lines() do
        linesData[lines] = line
        lines = lines + 1
    end

    reverse(linesData)

    data.paginator = paginator(page, 15, lines)
    data.logs = {}

    local current = 0
    local i = 0

    for index, line in ipairs(linesData) do
        if i >= data.paginator.offset then
            if current >= data.paginator.limit then
                break
            end
            table.insert(data.logs, json:unmarshal(line))
            current = current + 1
        end
        i = i + 1
    end

    http:render("logs.html", data)
end