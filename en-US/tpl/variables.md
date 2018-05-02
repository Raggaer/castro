---
name: Variables
---

# Syntax

Blocks must be wrapped between `{{` and `}}`.

## Variables

Template variables can be accessed and created. Creating a local variable follows the following syntax:

```html
{{ $var := 123 }}
```

Outputting local variables follows this syntax:

```html
<p>My name is {{ $var }}</p>
```

When you render a template you can optionally pass a lua table for some data. This lua table can be accessed on your template:

```lua
local data = {}

data.name = "Raggaer"

http:render("test.html", data)
```

```html
<h1>Welcome {{ .name }}</h1>
```

All these variables are always secure to use.

# Global variables

Variables that are available on each request:

- [widgets](#widgets)
- [microtime](#microtime)
- [session](#session)

## Widgets

Contains all the widgets html text. You can loop this variable to show the widgets:

```html
{{ range $index, $element := .widgets }}
    {{ $element }}
{{ end }}
```

## Microtime

Time that took to generate the page. This value is a string with the format `x seconds`.

**Microtime starts on a HTTP middleware so the value is always a little bit higher than what it should be**.

## Session

Provides access to the session map, usually you should access this by passing the session data to the template (with `http:render`) however sometimes you need to access the session map (for `i18n` for example)
