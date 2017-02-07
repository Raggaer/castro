local data = {}

data["validationError"] = session:getFlash("validationError")

http:render("search.html", data)