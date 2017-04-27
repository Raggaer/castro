---
Name: player
---

# Overview

Providess easy access to server players. You first need to create a new instance:

- [Player(name)](#new)
- [Player(id)](#new)

# new

Get a player by the given name

```lua
local data = Player("test")
```

This will return a new `player` metatable

# new

Get a player by the identifier

```lua
local data = Player(1)
```

This will return a new `player` metatable

# Player metatable

Provides access to the player information.

- [player:getAccountId()](#getaccountid)
- [player:isOnline()](#isonline)
- [player:getBankBalance()](#getbankbalance)
- [player:getStorageValue(key)](#getstoragevalue)
- [player:setStorageValue(key, value)](setstoragevalue)
- [player:getVocation()](#getvocation)
- [player:getTown()](#gettown)
- [player:getGender()](#getgender)
- [player:getLevel()](#getlevel)
- [player:getPremiumDays()](#getpremiumdays)
- [player:getName()](#getname)

# getAccountID

Returns the player account identifier.

```lua
local data = Player("test")
local accountId = data:getAccountId()
-- accountId = 13
```

# isOnline

Checks if the given player is online.

```lua
local data = Player("test")
local isOnline = data:isOnline()
-- isOnline = false
```

# getBankBalance

Returns the player bank balance.

```lua
local data = Player("test")
local balance = data:getBankBalance()
-- balance = 1000
```

# getStorageValue

Returns a storage value for the given player and key.

```lua
local data = Player("test")
local val = data:getStorageValue(1200)
-- val = 3000
```

# setStorageValue

Sets a storage value for the given player and key.

```lua
local data = Player("test")
data:setStorageValue(1200, 3000)
```

# getVocation

Returns the player vocation.

```lua
local data = Player("test")
local voc = data:getVocation()
--[[
voc.ID = 1
voc.Description = "A Mage"
voc.FromVoc = 0
voc.Name = "Mage"
]]--
```

The returned table contains these fields:

- ID: vocation identifier.
- Description: vocation description text.
- FromVoc: vocation FromVoc field.
- Name: vocation name.

All of these values are gathered from `vocations.xml` file.

# getTown

Returns the player town.

```lua
local data = Player("test")
local town = data:getTown()
--[[
town.ID = 1
town.Name = "Thais"
]]--
```

The returned table contains these fields:

- ID: town identifier.
- Name: town name.

# getGender

Returns the player gender.

```lua
local data = Player("test")
local gender = data:getGender()
-- gender = 1
```

# getLevel

Return the player current level.

```lua
local data = Player("test")
local level = data:getLevel()
-- level = 8
```

# getPremiumDays

Return the player current premium days.

```lua
local data = Player("test")
local gender = data:getPremiumDays()
-- gender = 365
```

# getName

Return the player name.

```lua
local data = Player(1)
local name = data:getName()
-- name = "test"
```