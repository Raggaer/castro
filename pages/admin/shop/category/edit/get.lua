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

    data.info = db:singleQuery("SELECT id, name, description FROM castro_shop_categories WHERE id = ?", http.getValues.id)

    if data.info == nil then
        http:redirect("/")
        return
    end

    data.success = session:getFlash("success")
    data.validationError = session:getFlash("validationError")
    data.id = data.info.id
    data.title = data.info.name
    data.text = data.info.description
    data.preview = data.info.description:parseBBCode()

    http:render("shopcategoryedit.html", data)
end