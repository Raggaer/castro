local data = {}

if http.getValues.town == nil then
    data.list = otbm:houseList()
else
    data.list = otbm:houseList(tonumber(http.getValues.town))
end

data.towns = otbm:townList()



http:render("houselist.html", data)

