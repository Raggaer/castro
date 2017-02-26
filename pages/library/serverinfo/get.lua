function get()
    local data = {}

    data.magic = config:get("rateMagic")
    data.skill = config:get("rateSkill")
    data.loot = config:get("rateLoot")
    data.stages = xml:unmarshalFile(serverPath .. "/data/XML/stages.xml")
    data.protection = config:get("protectionLevel")
    data.redskull = config:get("killsToRedSkull")

    http:render("serverinfo.html", data)
end