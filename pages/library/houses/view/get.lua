function get()
    local data = {}

    data.house = db:singleQuery([[
    SELECT c.name AS bidname, b.bid, b.bid_end, b.last_bid, a.name AS ownername, b.name, b.rent, b.size, b.beds, b.town_id, b.id FROM houses b 
    LEFT JOIN players a ON a.id = b.owner
    LEFT JOIN players c ON c.id = b.highest_bidder
    WHERE b.id = ?
    ]], http.getValues.id)

    if data.house == nil then
        http:redirect("/")
        return
    end

    data.town = otbm:townByID(tonumber(data.house.town_id))
    data.period = config:get("houseRentPeriod")
    data.logged = session:isLogged()
    data.error = session:getFlash("error")
    data.success = session:getFlash("success")

    if data.logged then
        data.characters = db:query("SELECT name FROM players WHERE account_id = ?", session:loggedAccount().ID)
    end

    http:render("viewhouse.html", data)
end

-- last_biD DINERO
-- hihest_bid player_id
-- ofertas deben ser mayores que bid