test = mysql:articleMultiple("SELECT id, title FROM articles")

print(config:getString("Motd"))
http:render("home.html")
