---
name: Map
---

# Map

Castro parses your Open Tibia binary map file at start-up. It then gets saved on the database `castro_map` table. Each map is saved by its name.

Parsing big maps (like realmap) can take longer on old setups.

If the map is already parsed Castro will check how old the map is and update if needed.

# Towns

All town data is gathered from the parsed map. If you add new towns please remove the entry from your `castro_map` table so a new map is generated.