function get()
    local data = {}

    data.magic = config:get("RateMagic")
    data.skill = config:get("RateSkill")
    data.loot = config:get("RateLoot")
    data.stages = xml:unmarshalFile(serverPath .. "/data/XML/stages.xml")
    data.protection = config:get("ProtectionLevel")
    data.redskull = config:get("RedSkullKills")

    http:render("serverinfo.html", data)
end