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