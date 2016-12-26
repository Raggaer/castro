test = mysql:articleMultiple("SELECT id, title FROM articles")
test[1].Title = "de que vas?"

test[1]:save("Title")

print(config:getString("Motd"))
http:render("home.html")
