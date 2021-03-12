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
    data.timeToDecreaseFrags = time:newDuration(config:get("timeToDecreaseFrags") * math.pow(10, 9))
    data.whiteSkullTime = time:newDuration(config:get("whiteSkullTime") * math.pow(10, 9))
    data.stages = config:get("experienceStages")

    if data.stages == nil then
        local stages = xml:unmarshalFile(serverPath .. "/data/XML/stages.xml")

        if stages.stages.config["-enabled"] == "1" then
            data.stages = {}
            for i, stage in ipairs(stages.stages.stage) do
                data.stages[i] = {
                    maxlevel = stage["-maxlevel"],
                    minlevel = stage["-minlevel"],
                    multiplier = stage["-multiplier"]
                }
            end
        end
    end

    http:render("serverinfo.html", data)
end
