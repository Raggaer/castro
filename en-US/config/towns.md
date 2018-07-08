---
name: Towns
---

# Towns

Manual list of the server towns. If you set `LoadMap` to true this field will have no use. Each element of the list contains the following fields:

- Name
- ID

Simple example on how the town list should look like inside your `config.toml` file:

```toml
[[Towns]]
  Name = "Thais"
  ID = 1

[[Towns]]
  Name = "Carlin"
  ID = 2

```