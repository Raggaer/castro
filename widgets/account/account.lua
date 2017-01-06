local data = {}

data["logged"] = session:isLogged()

if data["logged"] then
    data["account"] = session:loggedAccount()
end

return data
