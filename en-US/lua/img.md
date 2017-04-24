---
Name: img
---

# Img metatable

Provides access to image manipulation functions

- [img:new(width, height)](#new)

# new

Returns a new `goimage` instance. You can then manipulate this image as your needs

```lua
local test = img:new(500, 500)
```

# Goimage metatable

Provides access to image manipulation functions

- [goimage:writeText(text, color, size, x, y)](#writetext)
- [goimage:setBackground(filepath)](#setbackground)
- [goimage:save(path)](#save)

# writeText

Writes a text string to the given `goimage` instance. These are the list of mandatory parameters:

- text: text string.
- color: valid HEX color string. Example `#cccccc`.
- size: font size.
- x: X position on the image.
- y: Y position on the image.

```lua
local image = img:new(500, 500)

image:writeText("Hello World", "#D40000", 12, 40, 40)
```

`goimage` uses [go-font](https://blog.golang.org/go-fonts) as the default font family.

# setBackground

Sets the background for the given `goimage`

```lua
local image = img:new(500, 500)

img:setBackground("my_bg.jpg")
```

# save

Saves the given `goimage` result. Only `png` format is currently supported.

```lua
local image = img:new(500, 500)

image:writeText("Hello World", "#D40000", 12, 40, 40)
image:save("/images/example.png")
```