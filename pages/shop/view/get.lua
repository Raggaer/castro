function get()
    if not app.Shop.Enabled then
        http:redirect("/")
        return
    end
end