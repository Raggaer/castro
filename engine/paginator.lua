function paginator(page, perpage, count)
    local pg = {}
    pg.perpage = perpage
    pg.page = page
    pg.count = count
    pg.offset = page * perpage
    pg.limit = perpage
    pg.pages = {}
    if page <= 0 then
        pg.prev = false
    else
        pg.firstpage = newpage(0)
        pg.prev = true
        pg.prevnumber = page - 1
    end
    if (pg.page + 1) * pg.perpage >= pg.count then
        pg.last = false
    else
        local t = (count / perpage) - 1
        if math.floor(t) ~= t then
            t = math.floor(t) + 1
        end
        pg.lastpage = newpage(t)
        pg.lastnumber = page + 1
        pg.last = true
    end
    return pg
end

function newpage(num)
    local page = {}
    page.num = num
    return page
end