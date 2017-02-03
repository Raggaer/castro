if not session:isLogged() then
    return nil
end

local data = {}

data.admin = session:loggedAccount().castro.Admin

return data