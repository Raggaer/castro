function get()
    local data = {}

    data.list = xml:unmarshalFile(app.Main.Datapack .. "/data/spells/spells.xml")

    for i, spell in pairs(data.list.spells.instant) do
        if spell["-script"] == nil  then
            data.list.spells.instant[i] = nil
        else
            if string.find(spell["-script"], "monster/", 1) then
                data.list.spells.instant[i] = nil
            end
        end
    end

    http:render("spells.html", data)
end