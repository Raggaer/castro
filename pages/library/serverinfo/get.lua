function get()
    local data = {}

    data.experience = config:get("rateExp")
    data.magic = config:get("rateMagic")
    data.skill = config:get("rateSkill")
    data.loot = config:get("rateLoot")
    data.protection = config:get("protectionLevel")
    data.redskull = config:get("killsToRedSkull")
    data.blackskull = config:get("killsToBlackSkull")
    data.worldType = config:get("worldType")
    data.timeToDecreaseFrags = time:newDuration(config:get("timeToDecreaseFrags") * math.pow(10, 6))
    data.whiteSkullTime = time:newDuration(config:get("whiteSkullTime") * math.pow(10, 6))

    local stages = xml:unmarshalFile(serverPath .. "/data/XML/stages.xml")

    if stages.stages.config["-enabled"] == "1" then
        data.stages = stages
    end

    http:render("serverinfo.html", data)
end
