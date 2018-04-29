---
Name: guild
---

Provides easy access to server guilds. You first need to create a new instance:

- [Guild(name)](#player(name))
- [Guild(id)](#player(id))

# Guild(name)

Get a guild by the given name.

```lua
local data = Guild("test")
```

This will return a new `guild` metatable.

# Guild(id)

Get a guild by the unique identifier.

```lua
local data = Guild(1)
```

This will return a new `guild` metatable.

# Guild metatable

Provides access to the guild information.

- [guild:getOwner()](#getowner)

The table also contains some additional fields regarding guild information:

- `guild.ID`
- `guild.Name`
- `guild.Ownerid`
- `guild.Creationdata`
- `guild.Motd`

# getOwner

Retrieves the guild ownerid

```lua
local g = Guild("test")
local owner = g:getOwner()
-- owner = 10
```
