function get()
    if not session:isLogged() then
        http:redirect("/")
        return
    end

    local data = {}

    data.success = session:getFlash("success")
    data.towns = {}
    data.vocations = {}

    local vocations = xml:vocationList()
    local towns = otbm:townList()

    for _, town in pairs(towns) do
        for _, v in ipairs(app.Custom.ValidTownList) do
            if v == town.ID then
                table.insert(data.towns, town)
            end
        end
    end

    for _, voc in pairs(vocations) do
        for _, v in ipairs(app.Custom.ValidVocationList) do
            if v == voc.ID then
                table.insert(data.vocations, voc)
            end
        end
    end

    data.validationError = session:getFlash("validation-error")

    http:render("createcharacter.html", data)
end