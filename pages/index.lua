if _http_method == "POST" then
    render("home.html", _post)
else
    render("home.html", nil)
end