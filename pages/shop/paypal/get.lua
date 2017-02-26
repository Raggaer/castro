function get()
    local data = {}

    data.validationError = session:getFlash("validationError")

    http:render("paypal.html", data)
end