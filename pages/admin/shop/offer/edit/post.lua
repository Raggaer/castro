require "util"

function post()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end

    if not session:isAdmin() then
        http:redirect("/")
        return
    end

    local data = {}

    data.offer = db:singleQuery("SELECT id, category_id, price, description, name FROM castro_shop_offers WHERE id = ?", http.postValues.id)

    if data.offer == nil then
        http:redirect("/")
        return
    end

    data.category = db:singleQuery("SELECT id, name FROM castro_shop_categories WHERE id = ?", data.offer.category_id)

    if data.category == nil then
        http:redirect("/")
        return
    end

    if tonumber(http.postValues["offer-price"]) <= 0 then
        session:setFlash("validationError", "Invalid offer price range")
        http:redirect("/subtopic/admin/shop/offer/edit?id=" .. data.offer.id)
        return
    end

    if http.postValues["offer-name"] == "" then
        session:setFlash("validationError", "Offer name can not be empty")
        http:redirect("/subtopic/admin/shop/offer/edit?id=" .. data.offer.id)
        return
    end

    if http.postValues["offer-description"] == "" then
        session:setFlash("validationError", "Offer description can not be empty")
        http:redirect("/subtopic/admin/shop/offer/edit?id=" .. data.offer.id)
        return
    end

    http:parseMultiPartForm()

    local offerImage = http:formFile("offer-image")
    local offerImagePath = ""

    if offerImage ~= nil then

         if not offerImage:isValidPNG() then
            session:setFlash("validationError", "Offer image needs to be a valid png image")
            http:redirect("/subtopic/admin/shop/offer/edit?id=" .. data.offer.id)
            return
        end

        offerImage:saveFileAsPNG("public/images/offer-images/" .. http.postValues["offer-name"] .. ".png", 32, 32)
        offerImagePath = "/images/offer-images/" .. http.postValues["offer-name"] .. ".png"
    end

    db:execute(
        [[UPDATE castro_shop_offers
        SET description = ?,
        price = ?,
        name = ?,
        updated_at = ?,
        image = ?,
        give_item = ?,
        give_item_amount = ?,
        charges = ?,
        container_give_item = ?,
        container_give_amount = ?,
        container_give_charges = ?
        WHERE id = ?]],
        http.postValues["offer-description"],
        http.postValues["offer-price"],
        http.postValues["offer-name"],
        os.time(),
        offerImagePath,
        ternary(http.postValues["give-item"] == nil, 0, http.postValues["give-item"]),
        ternary(http.postValues["give-item-amount"] == nil, 0, http.postValues["give-item-amount"]),
        ternary(http.postValues["charges"] == nil, 0, http.postValues["charges"]),
        table.concat(explode(",", http.postValues["container-item[]"]), ","),
        table.concat(explode(",", http.postValues["container-item-amount[]"]), ","),
        table.concat(explode(",", http.postValues["container-item-charges[]"]), ","),
        data.offer.id
    )

    session:setFlash("success", "Shop offer edited")
    http:redirect("/subtopic/admin/shop/category?id=" .. data.category.id)
end
