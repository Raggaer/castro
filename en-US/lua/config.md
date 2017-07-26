---
name: config
---

# Config metatable

Provides access to the `config.lua` and `config.toml` files.

- [config:get(variable)](#get)
- [config:setCustom(key, data)](#setcustom)

# get

Retrieves a variable from the  `config.lua` file.

```lua
local url = config:get("url")
--- url = "www.google.com"
```

# setCustom

Sets a custom value on the `config.toml` file. You can pass a table, string, number or boolean as the data.

```lua
config:setCustom("test", "Hello World")
```

You can access the custom field of the `config.toml` file with the app variable, below is an example:

```lua
local test = app.Custom.test
--- test = "Hello World"
```

Setting a custom value will also be saved on the `config.toml` file.

# Config file access

In order to access `config.toml` file, you can access the `app` variable, this variable holds every field from the config file.

```lua
local mode = app.Mode
--- mode = "dev"
```

```lua
local paypal_enabled = app.PayPal.Enabled
--- paypal_enabled = true
```