require "util"

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

    data.offer = db:singleQuery("SELECT id, category_id, name, price, description, give_item, give_item_amount, container_item, container_give_item, container_give_amount FROM castro_shop_offers WHERE id = ?", http.getValues.id)
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

    if data.offer.container_give_item ~= "" and data.offer.container_give_item ~= nil then
        data.offer.containerItems = explode(",", data.offer.container_give_item)
    end

    if data.offer.container_give_amount ~= "" and data.offer.container_give_amount ~= nil then
        data.offer.containerAmounts = explode(",", data.offer.container_give_amount)
    end 

    http:render("editoffer.html", data)
end