require "bbcode"

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

    if http.postValues.submit and http.postValues.submit == "preview" then
        data.offer.description = http.postValues["offer-description"]
        data.parsedDescription = http.postValues["offer-description"]:parseBBCode()
        http:render("editoffer.html", data)
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

    if offerImage == nil then

        db:execute(
            "UPDATE castro_shop_offers SET description = ?, price = ?, name = ?, updated_at = ? WHERE id = ?",
            http.postValues["offer-description"],
            http.postValues["offer-price"],
            http.postValues["offer-name"],
            os.time(),
            data.offer.id
        )

    else

        if not offerImage:isValidPNG() then
            session:setFlash("validationError", "Offer image needs to be a valid png image")
            http:redirect("/subtopic/admin/shop/offer/edit?id=" .. data.offer.id)
            return
        end

        offerImage:saveFileAsPNG("public/images/offer-images/" .. http.postValues["offer-name"] .. ".png", 64, 64)

        db:execute(
            "UPDATE castro_shop_offers SET description = ?, price = ?, name = ?, updated_at = ?, image = ? WHERE id = ?",
            http.postValues["offer-description"],
            http.postValues["offer-price"],
            http.postValues["offer-name"],
            os.time(),
            "public/images/offer-images/" .. http.postValues["offer-name"] .. ".png",
            data.offer.id
        )

    end

    session:setFlash("success", "Shop offer edited")
    http:redirect("/subtopic/admin/shop/category?id=" .. data.category.id)
end