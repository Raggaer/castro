---
Name: outfit
---

# Outfit metatable

Provides access to the outfit metatable. Methods to generate outfits based on look values. Currently we ship with outfits up to `10.90` if you need to add a new outfit type head to `public/images/outfits/generator/` and add your outfit folder.

- [outfit:generate(looktype, lookfeet, looklegs, lookbody, looklegs, lookaddons)](#generate)

# generate 

Generates an outfit image with the given look values. This method returns a byte array as a string so you can use it to create image files.

```lua
local f = io.open("test.png", "w")
f:write(outfit:generate(128, 60, 30, 60, 30, 3))
f:close()
```