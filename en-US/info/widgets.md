---
name: Custom widgets
---

# Custom sidebar content

Castro uses lua to handle all your sidebar content. To create a new widget navigate to `widgets` directory and create a new folder with your desired name.

On your new folder you can then create a `widget_name.lua`.

Each widget must have a `function widget()` you can use the same metatables you would on a custom page with one addition: `widgets` metatable. You can call `widgets:render(template_name, data)` to render a widget template (just like you would render a page template)