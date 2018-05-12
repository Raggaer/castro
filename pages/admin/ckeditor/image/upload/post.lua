function post()
    if not session:isAdmin() then
        return
    end

    local file = http:formFile("upload")

    if file == nil then
        return
    end

    file:saveFile("public/ckeditor/images/" .. file.name)

    local data = {}
    data.uploaded = 1
    data.fileName = file.name
    data.url = "/ckeditor/images/" .. file.name

    http:write(json:marshal(data))
end