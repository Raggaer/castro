function post()
    if not session:isLogged() then
        http:redirect("/")
        return
    end

    local account = session:loggedAccount()

    local house = db:singleQuery([[
    SELECT c.name AS bidname, b.bid, b.bid_end, b.last_bid, a.name AS ownername, b.name, b.rent, b.size, b.beds, b.town_id, b.id FROM houses b 
    LEFT JOIN players a ON a.id = b.owner
    LEFT JOIN players c ON c.id = b.highest_bidder
    WHERE b.id = ?
    ]], http.postValues.id)

    if house == nil then
        http:redirect("/")
        return
    end

    if house.ownername ~= nil then
        http:redirect("/")
        return
    end

    local character = db:singleQuery("SELECT id, balance FROM players WHERE name = ? AND account_id = ?", url:decode(http.postValues.character), account.ID)

    if character == nil then
        http:redirect("/")
        return
    end

    if character.balance < tonumber(http.postValues.bid) then
        session:setFlash("error", "You need more gold coins to place that bid")
        http:redirect("/subtopic/library/houses/view?id=" .. house.id)
        return
    end

    if db:singleQuery("SELECT id FROM houses WHERE highest_bidder = ?", character.id) ~= nil then
        session:setFlash("error", "You already have a bid in place")
        http:redirect("/subtopic/library/houses/view?id=" .. house.id)
        return
    end

    -- There is no bid for the house
    if house.bidname == nil then
        db:execute(
            "UPDATE houses SET bid = ?, bid_end = ?, last_bid = ?, highest_bidder = ? WHERE id = ?",
            http.postValues.bid,
            os.time() + app.Custom.HouseBidOpenTime,
            http.postValues.bid,
            character.id,
            house.id
        )

        session:setFlash("success", "Bid placed")
        http:redirect("/subtopic/library/houses/view?id=" .. house.id)

        return
    end

    if house.highest_bidder == character.id then
        session:setFlash("error", "You cannot outbid your own bid")
        http:redirect("/subtopic/library/houses/view?id=" .. house.id)
        return
    end

    if tonumber(http.postValues.bid) < house.last_bid then
        session:setFlash("error", "Your bid needs to be higher than the last bid. The last bid is currently " .. house.last_bid)
        http:redirect("/subtopic/library/houses/view?id=" .. house.id)
        return
    end

    db:execute(
        "UPDATE houses SET bid = bid + ?, last_bid = ?, highest_bidder = ? WHERE id = ?",
        http.postValues.bid,
        http.postValues.bid,
        character.id,
        house.id
    )
end
