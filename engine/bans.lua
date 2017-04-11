-- Castro needs to be restarted for changes in this file to take effect
local config = {}

config.quickList = {
    --[[
    You may add or modify the quick bans drop down options here.

    duration should be a string ("4h") or number (seconds)

    Valid types are:
    namelock: namelocks the player
    account_ban: ban the player's account
    ip_ban: ban the last IP player connected from
    account_ip_ban: ban account and last known IP
    ]]
    {id = 1, reason = "Offensive Name", type = "namelock"},
    {id = 2, reason = "Unsuitable Name", type = "namelock"},
    {id = 3, reason = "Name Inciting Rule Violation", type = "namelock"},
    {id = 4, reason = "Offensive Statement", type = "account_ban", duration = "4h"},
    {id = 5, reason = "Spamming", type = "account_ban", duration = "24h"},
    {id = 6, reason = "Illegal Advertising", type = "account_ip_ban", duration = "48h"},
    {id = 7, reason = "Inciting Rule Violation", type = "account_ban", duration = "24h"},
    {id = 8, reason = "Bug Abuse", type = "account_ip_ban", duration = "72h"},
    {id = 9, reason = "Using Unofficial Software to Play", type = "account_ban", duration = "24h"},
    {id = 10, reason = "Multi-Clienting", type = "account_ip_ban", duration = "24h"},
    {id = 11, reason = "Account Trading or Sharing", type="account_ban", duration = "24h"},
    {id = 12, reason = "Destructive Behaviour", type = "account_ip_ban", duration = "48h"},
    {id = 13, reason = "Excessive Unjustified Player Killing", type = "account_ban", duration = "48h"},
}

return config
