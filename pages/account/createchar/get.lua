function get()
    if not session:isLogged() then
        http:redirect("/")
        return
    end

    local data = {}

    data.success = session:getFlash("success")
    data.vocations = xml:vocationList(true)
    data.towns = {}

    local towns = otbm:townList()

    for _, town in pairs(towns) do
        for _, v in ipairs(app.Custom.ValidTownList) do
            if v == town.ID then
                table.insert(data.towns, town)
            end
        end
    end

    data.validationError = session:getFlash("validation-error")

    http:render("createcharacter.html", data)
end