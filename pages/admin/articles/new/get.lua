function get()
-- Block access for anyone who is not admin
if not session:isLogged() or not session:isAdmin() then
	http:redirect("/")
	return
end

local data = {}
data.validationError = session:getFlash("validationError")

http:render("newarticle.html", data)
end