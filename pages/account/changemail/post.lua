if not session:isLogged() then
    http:redirect("/")
    return
end

local account = session:loggedAccount()

if http.postValues["account-email"] ~= account.Email then
    session:setFlash("validationError", "Wrong account email")
    http:redirect("/subtopic/account/changemail")
    return
end

db:execute("UPDATE accounts SET email = ? WHERE id = ?", http.postValues["new-email"], account.ID)
session:setFlash("success", "Email changed")
http:redirect("/subtopic/account/dashboard")
