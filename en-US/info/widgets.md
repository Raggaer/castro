---
name: Custom widgets
---

# Custom sidebar content

Castro uses lua to handle all your sidebar content. To create a new widget navigate to `widgets` directory and create a new folder with your desired name.

On your new folder you can then create a `widget_name.lua`.

Each widget must have a `function widget()` you can use the same metatables you would on a custom page with one addition: `widgets` metatable. You can call `widgets:render(template_name, data)` to render a widget template (just like you would render a page template). Below is an example on how a simple widget will look like:

```lua
function widget()
    local data = {}

    data.top = db:query("SELECT name, level, group_id FROM players WHERE group_id = 1 ORDER BY level DESC LIMIT 5")

    widgets:render("toplevel.html", data)
end
```

In this example we load the top five players from database and pass it to our template.

```html
<div class="info-panel black">
    <h5>Top players</h5>
    <hr class="white-hr">
    <ul>
        {{ if .top }}
            {{ range $index, $element := .top }}
                <li><a class="light" href="{{ url "subtopic" "community" "view" }}?name={{ urlEncode $element.name }}">{{ $element.name }}</a> ({{ $element.level }})</li>
            {{ end }}
        {{ else }}
            No players made
        {{ end }}
    </ul>
</div>
```

We check if our top variable is not nil, loop over that variable and show the players