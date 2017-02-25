function get()
    local data = {}

    data.magic = config:get("RateMagic")
    data.skill = config:get("RateSkill")
    data.loot = config:get("RateLoot")
    data.stages = xml:unmarshalFile(serverPath .. "/data/XML/stages.xml")

    http:render("serverinfo.html", data)
end