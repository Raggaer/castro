function paginator(page, perpage, count)
    local pg = {}
    pg.perpage = perpage
    pg.page = page
    pg.count = count
    pg.limit = page * perpage
    pg.offset = perpage
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
        pg.lastpage = newpage((count / perpage) - 1)
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