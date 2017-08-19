---
Name: map
---

# Map metatable

Provides access to map related functions. All the methods will use your current map used by your `config.lua` file.

- [otbm:encode()](#encode)

# encode

Encodes the current running map. Returns a byte array with the needed information for Castro. The map is saved on database over the `castro_map` table.

```lua
otbm:encode()
```

This is can be big resource consuming task (depending on your map). Use at your own risk.