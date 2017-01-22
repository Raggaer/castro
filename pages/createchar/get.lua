if not session:isLogged() then
    http:redirect("/")
    return
end

local data = {}

data.success = session:getFlash("success")
data.vocations = xml:vocationList(true)
data.towns = otbm:townList()
data.validationError = session:getFlash("validation-error")

http:render("createcharacter.html", data)