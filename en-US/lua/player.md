---
Name: player
---

# Overview

Provides easy access to server players. You first need to create a new instance:

- [Player(name)](#player(name))
- [Player(id)](#player(id))

# Player(name)

Get a player by the given name.

```lua
local data = Player("test")
```

This will return a new `player` metatable.

# Player(id)

Get a player by the unique identifier.

```lua
local data = Player(1)
```

This will return a new `player` metatable.

# Player metatable

Provides access to the player information.

- [player:getGuild()](#getguild)
- [player:getAccountId()](#getaccountid)
- [player:isOnline()](#isonline)
- [player:getBankBalance()](#getbankbalance)
- [player:setBankBalance()](#setbankbalance)
- [player:getStorageValue(key)](#getstoragevalue)
- [player:setStorageValue(key, value)](setstoragevalue)
- [player:getVocation()](#getvocation)
- [player:getTown()](#gettown)
- [player:getGender()](#getgender)
- [player:getLevel()](#getlevel)
- [player:getPremiumDays()](#getpremiumdays)
- [player:getName()](#getname)
- [player:getExperience()](#getexperience)
- [player:getCapacity()](#getcapacity)
- [player:getCustomField()](#getcustomfield)
- [player:setCustomField()](#setcustomfield)

The table also contains some additional fields regarding player information:

- `player.ID`
- `player.Name`
- `player.Level`
- `player.Sex`
- `player.Vocation`
- `player.Town_id`
- `player.Account_id`
- `player.Experience`

# getGuild

Returns the player guild id.

```lua
local data = Player("test")
local g = data:getGuild()
-- g = 10
```

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

# setBankBalance

Updates the player bank balance.

```lua
local data = Player("test")
local balance = data:setBankBalance(100)
-- balance = 100
```

This method overwrites the bank balance. If you want to augment the player balance you can do the following:

```lua
local data = Player("test")
local balance = data:setBankBalance(data:getBankBalance() + 100)
-- balance = 1100
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

# getExperience

Return the player current experience value.

```lua
local data = Player(1)
local experience = data:getExperience()
-- experience = 4800
```

# getCapacity

Return the player current capacity value.

```lua
local data = Player(1)
local cap = data:getCapacity()
-- capacity = 100
```

# getCustomField

Retrieves the given field as a string.

```lua
local data = Player(1)
local points = tonumber(data:getCustomField("points"))
-- points = 100
```

# setCustomField

Sets the given field on the player table

```lua
local data = Player("Test")
data:setCustomField("level", 999)
```