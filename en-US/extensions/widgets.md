---
name: Widgets
---

# Widgets

Widgets in an extension should be placed in a folder named `widgets`. Folders inside `widgets` will be added to the widget application list.
You can add unlimited amount of widgets for an extension by simply adding more folders inside the `widgets` main folder.

## Files

Widgets consists of 2 types of files, html for templates and Lua for logic. The Lua and html files for extension pages works the same way as [custom widgets](/docs/info/widgets).

### Folder structure

```html
my_extension /
    / widgets
        / test
            test.lua
            my-widget-template.html
```

## Overriding defaults

It is not possible to override widgets using an extension.