require "extensionhooks"

function post()
    if not session:isLogged() then
        http:redirect("/subtopic/login")
        return
    end

    -- Extension hook
    executeHook("onLogout", session:loggedAccount())

    session:destroy()
    session:setFlash("success", "Logged out")

    http:redirect("/subtopic/login")
end