-- Reverse the given table
function reverse(tbl)
    for i=1, math.floor(#tbl / 2) do
        tbl[i], tbl[#tbl - i + 1] = tbl[#tbl - i + 1], tbl[i]
    end
end

-- Return application running absolute URL
function runningURL()
    local mode = "http"

    if app.SSL.Enabled then
        mode = "https"
    end

    return string.format("%s://%s:%s", mode, app.URL, app.Port)
end