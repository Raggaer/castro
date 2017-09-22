function post() 
    if not session:isAdmin() then
        return
    end
    
    local file = http:formFile("file")

    if file == nil then
        return
    end

    file:saveFile("public/froala/files/" .. file.name)
    http:write(json:marshal({ link = "/froala/files/" .. file.name }))
end