if not session:isLogged() then
    http:redirect("/")
    return
end

local data = {}

data.vocations = xml:vocationList(true)
data.towns = otbm:townList()

http:render("createcharacter.html", data)