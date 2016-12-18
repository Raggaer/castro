local test = {}
test["hola"] = 120
test["more"] = {}
test["more"]["name"] = "alvaro"
test["more"][0] = "hola"

renderTemplate("home.html", test)