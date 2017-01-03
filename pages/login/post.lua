account = db:query("SELECT name, email FROM accounts WHERE name = ? AND password = ?", http.postValues["account-name"], crypto:sha1(http.postValues.password))

if account == nil then
    session:setFlash("validationError", "Wrong account name or password")
    http:redirect("/subtopic/login")
end

