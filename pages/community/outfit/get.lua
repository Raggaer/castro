function get()
    local looktype = tonumber(http.getValues.looktype)
    local lookfeet = tonumber(http.getValues.lookfeet)
    local looklegs = tonumber(http.getValues.looklegs)
    local lookbody = tonumber(http.getValues.lookbody)
    local lookhead = tonumber(http.getValues.lookhead)
    local lookaddons = tonumber(http.getValues.lookaddons)

    if looktype == nil or (looktype < 0 and looktype > 132) then
        return
    end

    if lookfeet == nil or (lookfeet < 0 and lookfeet > 132) then
        return
    end

    if looklegs == nil or (looklegs < 0 and looklegs > 132) then
        return
    end

    if lookbody == nil or (lookbody < 0 and lookbody > 132) then
        return
    end

    if lookhead == nil or (lookhead < 0 and lookhead > 132) then
        return
    end

    if lookaddons == nil or (lookaddons < 1 and lookaddons > 3) then
        return
    end

    local outfitpath = string.format("public/images/outfits/%d_%d_%d_%d_%d_%d.png", looktype, lookfeet, looklegs, lookbody, lookhead, lookaddons)

    if file:exists(outfitpath) then
        http:serveFile(outfitpath)
        return
    end

    local out = outfit:generate(
        looktype,
        lookfeet,
        looklegs,
        lookbody,
        lookhead,
        lookaddons
    )

    local outfile = io.open(outfitpath, "w")

    outfile:write(out)
    outfile:close()
    http:serveFile(outfitpath)
end
