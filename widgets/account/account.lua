local data = {}

data["logged"] = session:isLogged()
data["account"] = db:query("SELECT points FROM castro_accounts WHERE name = ?", session:get("logged-account"))

return data
