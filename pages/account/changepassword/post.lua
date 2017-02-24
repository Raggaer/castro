function post()

if not session:isLogged() then
    http:redirect("/")
    return
end

local account = session:loggedAccount()

if crypto:sha1(http.postValues["account-password"]) ~= account.Password then
    session:setFlash("validationError", "Wrong account password")
    http:redirect("/subtopic/account/changepassword")
    return
end

db:execute("UPDATE accounts SET password = ? WHERE id = ?", crypto:sha1(http.postValues["new-password"]), account.ID)
session:setFlash("success", "Password changed")
http:redirect("/subtopic/account/dashboard")

    end