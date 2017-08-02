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

    data.category = db:singleQuery("SELECT id, name FROM castro_shop_categories WHERE id = ?", http.postValues.id)

    if data.category == nil then
        http:redirect("/")
        return
    end

    if http.postValues.submit and http.postValues.submit == "preview" then
        data.description = http.postValues["offer-description"]
        data.parsedDescription = http.postValues["offer-description"]:parseBBCode()
        http:render("newoffer.html", data)
        return
    end

    if http.postValues["offer-price"] == nil then
        session:setFlash("validationError", "Invalid offer price range")
        http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
        return
    end

    if tonumber(http.postValues["offer-price"]) ~= nil and tonumber(http.postValues["offer-price"]) <= 0 then
        session:setFlash("validationError", "Invalid offer price range")
        http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
        return
    end

    if http.postValues["offer-name"] == "" then
        session:setFlash("validationError", "Offer name can not be empty")
        http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
        return
    end

    if http.postValues["offer-description"] == "" then
        session:setFlash("validationError", "Offer description can not be empty")
        http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
        return
    end

    if db:singleQuery("SELECT 1 FROM castro_shop_offers WHERE name = ?", http.postValues["offer-name"]) ~= nil then
        session:setFlash("validationError", "Offer name is already in use")
        http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
        return
    end

    http:parseMultiPartForm()

    local offerImage = http:formFile("offer-image")

    if offerImage ~= nil then

        if not offerImage:isValidPNG() then
            session:setFlash("validationError", "Offer image can only be .png")
            http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
            return
        end

        offerImage:saveFileAsPNG("public/images/offer-images/" .. http.postValues["offer-name"] .. ".png", 64, 64)

        db:execute(
            "INSERT INTO castro_shop_offers (category_id, description, price, name, created_at, updated_at, image) VALUES (?, ?, ?, ?, ?, ?, ?)",
            data.category.id,
            http.postValues["offer-description"],
            http.postValues["offer-price"],
            http.postValues["offer-name"],
            os.time(),
            os.time(),
            "public/images/offer-images/" .. http.postValues["offer-name"] .. ".png"
        )

    else

        db:execute(
            "INSERT INTO castro_shop_offers (category_id, description, price, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
            data.category.id,
            http.postValues["offer-description"],
            http.postValues["offer-price"],
            http.postValues["offer-name"],
            os.time(),
            os.time()
        )
    end

    session:setFlash("success", "Shop offer created")
    http:redirect("/subtopic/admin/shop/category?id=" .. data.category.id)
end