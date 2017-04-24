---
Name: http
---

# Http metatable

Provides access to HTTP related functions.

- [http:redirect(url, header)](#redirect)
- [http:render(template, data)](#render)

# redirect

Redirects the user to the given location. You can provide an optional header. By default all redirects are done using a `302` header.

```lua
http:redirect("/subtopic/test")
```

Redirecting an user does not stop the execution of the page. You must return on each redirect.

# render

Renders the given template. You can pass a lua table for the template data. This data is HTML escaped by the templating engine.

```lua
local data = {}

data.name = "Raggaer"

http:render("test.html", data)
```