---
Name: image
---

# Image metatable

Provides access to image manipulation functions.

- [image:new(width, height)](#new)

# new

Returns a new `goimage` instance. You can then manipulate this image.

```lua
local test = img:new(500, 500)
```

# Goimage metatable

Provides access to image manipulation functions:

- [goimage:encode()](#encode)
- [goimage:writeText(text, color, size, x, y, optional font)](#writetext)
- [goimage:setBackground(filepath)](#setbackground)
- [goimage:save(path)](#save)

# encode

Encodes the image into a byte array string. You can later use this string:

```lua
local test = image:new(500, 500)
test:writeText("Hello World", "#D40000", 12, 40, 40)

local v = test:encode()

http:write(v)
```

On the example we create a image and serve it on the fly (without saving it as a file)

# writeText

Writes a text string to the given `goimage` instance. These are the list of mandatory parameters:

- text: text string.
- color: valid HEX color string. Example `#cccccc`.
- size: font size.
- x: X position on the image.
- y: Y position on the image.

```lua
local image = image:new(500, 500)

image:writeText("Hello World", "#D40000", 12, 40, 40)
```

`goimage` uses [go-font](https://blog.golang.org/go-fonts) as the default font family.

You can however specify your custom font file path:

```lua
local image = image:new(500, 500)

image:writeText("Hello World", "#D40000", 12, 40, 40, "Martel.ttf")
```

# setBackground

Sets the background for the given `goimage`

```lua
local image = image:new(500, 500)

image:setBackground("my_bg.jpg")
```

# save

Saves the given `goimage` result. Only `png` format is currently supported.

```lua
local image = image:new(500, 500)

image:writeText("Hello World", "#D40000", 12, 40, 40)
image:save("/images/example.png")
```
