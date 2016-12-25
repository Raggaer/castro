test = mysql:articleSingle("id = ?", 1)
test.Title = "Hello11211"
test:save()
print(config:getString("Motd"))
http:render("home.html")
