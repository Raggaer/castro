require "bbcode"

function get()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    local data = {}

    data.categories = db:query("SELECT id, name, description FROM castro_shop_categories ORDER BY id")

    for _, category in ipairs(data.categories) do
        category.parsedDescription = category.description:parseBBCode()
        category.offers = db:query("SELECT id, image, name, description, price FROM castro_shop_offers WHERE category_id = ?", category.id)
        for _, offer in ipairs(category.offers) do
            offer.parsedDescription = offer.description:parseBBCode()
        end
    end

    data.success = session:getFlash("success")

    http:render("shopview.html", data)
end