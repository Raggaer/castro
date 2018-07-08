---
name: Config
---

Below are the main configuration fields:

- [Mode](#mode)
- [CheckUpdates](#checkupdates)
- [Port](#port)
- [URL](#url)
- [Datapack](#datapack)
- [LoadMap](#loadmap)
- [MapHouseFile](#maphousefile)

# Mode

Sets the castro running mode. Possible values are

- `prod`
- `dev`

While on `dev` mode Castro will reload all pages, widgets and config file on each request. Dont run Castro using `dev` mode while your site is public available, the `dev` mode has a big performance and memory hit on your system and should only be used for local development.

# CheckUpdates

If `true` checks how many commits behind you are running Castro at start-up.

# Port

Determines the port where the HTTP server should listen to. By default browsers will look for port `80`.

# URL

Sets the URL for all links.

# Datapack

Path of your server datapack. Where `config.lua` is located. This path is set at the installation process.

# LoadMap

If enabled your map data will be loaded from your `*.otbm` file, sometimes this wont work (using old map protocol for example), in this case you might want to disable this and manually configure towns and house map file

# MapHouseFile

Name of your map house file, this field is only needed when setting the field `LoadMap` to false