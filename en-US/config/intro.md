---
name: Config
---

Castro uses a configuration file to handle all systems. The file `config.toml` is generated at the installation process, some values like for example the cookie hash key are also generated at the installation process.

You can access any of these fields from your `lua` files using the  `app` variable:

```lua
local mode = app.Mode
--- mode = "dev"
```

Each subtopic, widget or extension can also have its own configuration file, simply create a `config.lua` inside the directory, this way you can overwrite config values while keeping the main file clean, you can even use all the metatables castro exposes.

For example, on the create character subtopic there is a  `config.lua` file that sets the valid vocations and towns for new characters:

```lua
-- Valid town id list
app.Custom.ValidTownList = {1, 2, 3, 4}

-- Valid vocation id list
app.Custom.ValidVocationList = {1, 2, 3, 4}
```

Check the rest of the pages inside the config topic to learn more about each value of the `config.toml` file.