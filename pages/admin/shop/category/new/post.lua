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

    if http.postValues.action ~= "add" then
        http:redirect("/")
        return
    end

    if http.postValues.title:len() > 45 then
        session:setFlash("validationError", "Category title must have less than 45 characters")
        http:redirect("/subtopic/admin/shop/category/new")
        return
    end

    if http.postValues.text:len() > 255 then
        session:setFlash("validationError", "Category description must have less than 255 characters")
        http:redirect("/subtopic/admin/shop/category/new")
        return
    end

    db:execute("INSERT INTO castro_shop_categories (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)", http.postValues.title, http.postValues.text, os.time(), os.time())
    session:setFlash("success", "Category created")
    http:redirect("/subtopic/admin/shop")
end