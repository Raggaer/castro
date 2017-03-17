-- Return application running absolute URL
function runningURL()
    local mode = "http"

    if app.SSL.Enabled then
        mode = "https"
    end

    return string.format("%s://%s:%s", mode, app.URL, app.Port)
end