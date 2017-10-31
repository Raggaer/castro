---
name: Map
---

# Map

Castro parses your Open Tibia binary map file at start-up. It then gets saved on the database `castro_map` table. Each map is saved by its name.

Parsing big maps (like realmap) can take longer on old setups.

If the map is already parsed Castro will check how old the map is and update if needed.

## Towns

All town data is gathered from the parsed map. If you add new towns you have two options to reload the map:

- Restart: Castro will check for any updates on your map file.
- Wait: Castro service will eventually notice about your map beeing updated.

## Reload

While Castro is running, as you have seen, the map is reloaded when you make changes to it. The period of time between changes is configured on the [map watcher configuration field](/docs/config/mapwatch)