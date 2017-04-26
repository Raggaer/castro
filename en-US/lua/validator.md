---
name: validator
---

# Validator metatable

Provides access to data validation functions.

- [validator:validQRToken(token, secret)](#validqrtoken)
- [validator:validGuildName(name)](#validguildname)
- [validator:validGuildRank(rank)](#validguildrank)
- [validator:validVocation(name or id, base = false)](#validvocation)
- [validator:validTown(name or id)](#validtown)
- [validator:validUsername(name)](#validusername)

# validQRToken

Validates the given token against a secret key. Castro uses this function for the two-factor authentication.

```lua
local status = validator:validQRToken("token", "secret")
-- status = false
```

# validGuildName

Validates the given guild name. The name must be between 5 to 20 characters long and must compile against this regular expression: 

```
^[a-zA-Z ]+$
```

```lua
local valid = validator:validGuildName("This is not valid!!")
-- valid = false
```

# validGuildRank

Validates the given guild rank. The rank must be between 5 to 15 characters long and must compile against this regular expression: 

```
^[a-zA-Z ]+$
```

```lua
local valid = validator:validGuildRank("This is not valid!!")
-- valid = false
```

# validVocation

Validates the given vocation by its identifier or name. You can also pass a second parameter to only validate base vocations.
 
```lua
local valid = validator:validVocation("Sorcerer")
-- valid = true
``` 

```lua
local valid = validator:validVocation(10)
-- valid = true
```

# validTown

Validates the given town by its identifier or name.

```lua
local valid = validator:validTown("Thais")
-- valid = true
``` 

```lua
local valid = validator:validTown(10)
-- valid = true
```

# validUsername

Validates the given character name. Must compile against this regular expression: 

```
^[a-zA-Z ]+$
```

```lua
local valid = validator:validUsername("Not_valid!")
-- valid = false
```

