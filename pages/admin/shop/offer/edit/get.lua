require "bbcode"

function get()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isAdmin() then
        http:redirect("/")
        return
    end

    local data = {}

    data.offer = db:singleQuery("SELECT id, category_id, name, price, description FROM castro_shop_offers WHERE id = ?", http.getValues.id)

    data.validationError = session:getFlash("validationError")

    if data.offer == nil then
        http:redirect("/")
        return
    end

    data.category = db:singleQuery("SELECT name FROM castro_shop_categories WHERE id = ?", data.offer.category_id)

    if data.category == nil then
        http:redirect("/")
        return
    end

    data.parsedDescription = data.offer.description:parseBBCode()

    http:render("editoffer.html", data)
end