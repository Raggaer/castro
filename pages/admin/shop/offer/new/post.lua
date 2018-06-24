require "bbcode"
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

    data.category = db:singleQuery("SELECT id, name FROM castro_shop_categories WHERE id = ?", http.postValues.id)

    if data.category == nil then
        http:redirect("/")
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
    local offerImagePath = ""

    if offerImage ~= nil then

        local validImage = offerImage:isValidExtension("image/png") or offerImage:isValidExtension("image/jpeg") or offerImage:isValidExtension("image/gif")

        if not validImage then
            session:setFlash("validationError", "Offer image can only be .png, .jpeg or .gif")
            http:redirect("/subtopic/admin/shop/offer/new?categoryId=" .. data.category.id)
            return
        end

        offerImage:saveFileAsPNG("public/images/offer-images/" .. http.postValues["offer-name"] .. ".png", 32, 32)
        offerImagePath = "public/images/offer-images/" .. http.postValues["offer-name"] .. ".png"
    end

    db:execute(
        [[INSERT INTO castro_shop_offers
        (category_id, description, price, name, created_at, updated_at, image, give_item, give_item_amount, charges, container_give_item, container_give_amount, container_give_charges)
        VALUES
        (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)]],
        data.category.id,
        http.postValues["offer-description"],
        http.postValues["offer-price"],
        http.postValues["offer-name"],
        os.time(),
        os.time(),
        offerImagePath,
        ternary(http.postValues["give-item"] == nil, 0, http.postValues["give-item"]),
        ternary(http.postValues["give-item-amount"] == nil, 0, http.postValues["give-item-amount"]),
        ternary(http.postValues["charges"] == nil, 0, http.postValues["charges"]),
        ternary(http.postValues["container-item[]"] == "", "", table.concat(explode(",", http.postValues["container-item[]"]), ",")),
        ternary(http.postValues["container-item-amount[]"] == "", "", table.concat(explode(",", http.postValues["container-item-amount[]"]), ",")),
        ternary(http.postValues["container-item-charges[]"] == "", "", table.concat(explode(",", http.postValues["container-item-charges[]"]), ","))
    )

    session:setFlash("success", "Shop offer created")
    http:redirect("/subtopic/admin/shop/category?id=" .. data.category.id)
end
