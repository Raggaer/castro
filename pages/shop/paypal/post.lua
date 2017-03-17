function post()
    if not app.PayPal.Enabled then
        http:redirect("/")
        return
    end

    -- print(paypal:createPayment(paypalList["Test Package"].name, "test", paypalList["Test Package"].price))
end