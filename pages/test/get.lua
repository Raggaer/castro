function get()
    local data = {}
    data.name = "alvaro"
    data.second = {}
    data.second.name = "jaime"
    data.test = {}
    data.test[0] = 12
    data.test[1] = 13

    http:render("test.html", data)
end