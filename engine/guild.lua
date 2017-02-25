function isGuildOwner(accountid, guild)
    local characters = db:query("SELECT id FROM players WHERE account_id = ?", accountid)

    for _, val in pairs(characters) do
        if val.id == tonumber(guild.ownerid) then
            return true
        end
    end

    return false
end