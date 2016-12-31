test = db:query("SELECT title FROM articles WHERE id = ?", 2)
print(test[1].title)
print(config:getString("Motd"))
http:render("home.html")
