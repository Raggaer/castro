function get()
    if not session:isAdmin() then
        return
    end

    local list = {}

    list.list = {}

    for i, v in pairs(file:getFiles("public/froala/images")) do
        table.insert(list.list, {
            name = v,
            id = i,
            url = "/froala/images/" .. v,
            type = "image",
            thumb = "/froala/images/" .. v
        })
    end

    http:write(json:marshal(list):sub(1, -2):sub(9))
end