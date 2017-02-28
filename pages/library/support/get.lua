function get()
    local data = {}
    local groups = xml:unmarshalFile(serverPath .. "/data/XML/groups.xml")

    data.list = {}

    for _, group in pairs(groups.groups.group) do
        if tonumber(group["-access"]) >= 1 then
            data.list[group["-name"]] = db:query("SELECT name, lastlogin FROM players WHERE group_id = ?", group["-id"])
            if data.list[group["-name"]] then
                for _, player in pairs(data.list[group["-name"]]) do
                    player.lastlogin = time:parseUnix(player.lastlogin)
                end
            end
        end
    end

    http:render("support.html", data)
end