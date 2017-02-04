local info = {}

info["magic"] = config:get("RateMagic")
info["loot"] = config:get("RateLoot")
info["online"] = db:query("SELECT COUNT(*) as online FROM players_online", true).online

return info, true