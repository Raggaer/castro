function get()
    local data = {}
    data.name = "alvaro"

    local chart = {
        labels = {"A", "B", "C"},
        datasets = {
            {
                label = 'Players online',
                data = {1, 2, 3},
                backgroundColor = 'rgba(0, 140, 186, 0.2)',
                borderColor = 'rgba(0, 140, 186, 1)',
                borderWidth = 1,
            }
        }
    }
    print(json:marshal(chart))


    http:render("test.html", data)
end