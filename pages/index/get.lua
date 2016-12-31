articles = db:query("SELECT title, text, created_at FROM articles ORDER BY id DESC LIMIT 5")
http:render("home.html", { list = articles })
