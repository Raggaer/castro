function paygolPackage(name, type, price, points)
    local pkg = {}

    pkg.name = name
    pkg.type = type
    pkg.price = price
    pkg.points = points

    return pkg
end

-- PayGol package list
paygolList = {}

-- Package name, payment method, price in currency, points given
paygolList["Test Package"] = paygolPackage("Test Package", "sms", 12, 10)
