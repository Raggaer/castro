---
name: Loops
---

# Loops

It is possible to loop variables inside a template. The syntax is the following:

```html
{{ range $index, $element := .var }}

    <h1>{{ $element }}</h1>

{{ end }}
```

The `$index` variable is optional.

When inside a loop you cant access global variables like you would normally do using `{{ .var }}`. Instead you use the following syntax `{{ $.var }}`:

```lua
local data = {}

data.name = "Raggaer"
data.text = "This is a variable"
data.list[1] = "First"
data.list[2] = "Second"

http:render("test.html", data)
```

```html
<h1>Welcome {{ .name }}</h1>
{{ range $index, $element := .list }}

    <p>{{ $element }} and {{ $.text }}</p>

{{ end }}
```