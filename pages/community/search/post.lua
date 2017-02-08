if http.postValues["character-name"] == "" then
    session:setFlash("validationError", "Search query cant be empty")
    http:redirect("/subtopic/community/search")
    return
end

local data = {}

data.list = db:query("SELECT name FROM players WHERE name LIKE ?", "%" .. http.postValues["character-name"] .. "%", true)

if data.list == nil then
    session:setFlash("validationError", "No results found")
    http:redirect("/subtopic/community/search")
    return
end

data.query = http.postValues["character-name"]

http:render("search.html", data)