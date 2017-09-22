function post()
    if not session:isAdmin() then
        return
    end
    
    local file = http:formFile("file")

    if file == nil then
        return
    end

    if not file:isValidPNG() then
        http:write(json:marshal({ error = "Only PNG images are allowed"}))
        return
    end

    file:saveFileAsPNG("public/froala/images/" .. file.name)
    http:write(json:marshal({ link = "/froala/images/" .. file.name }))
end
