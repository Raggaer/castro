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
        http:redirect("/subtopic/admin/shop/category/edit?id=" .. http.postValues.id)
        return
    end

    if http.postValues.text:len() > 255 then
        session:setFlash("validationError", "Category description must have less than 255 characters")
        http:redirect("/subtopic/admin/shop/category/edit?id=" .. http.postValues.id)
        return
    end

    db:execute("UPDATE castro_shop_categories SET name = ?, description = ?, updated_at = ? WHERE id = ?", http.postValues.title, http.postValues.text, os.time(), http.postValues.id)
    session:setFlash("success", "Category updated")
    http:redirect("/subtopic/admin/shop")
end