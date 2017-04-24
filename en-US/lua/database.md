---
Name: database
---

# Database metatable

Provides access to the server database. Using prepared statements

* [db:singleQuery(query, args, cache)](#singlequery)
* [db:query(query, args, cache)](#query)
* [db:execute(query)](#execute)

# singleQuery

Excecutes a query returning only one result. The cache value is optional
Each query parameter should be a `?` symbol. Then passing the real value as an argument.

```lua
local name = "test"
local articles = db:singleQuery("SELECT id FROM articles WHERE name = ?", name, false)
--[[ articles.id = 1 ]]--
```

singleQuery also returns the cache status. If you need to check if the query is returned from the cache or not

```lua
local name = "test"
local articles, cache = db:singleQuery("SELECT id FROM articles WHERE name = ?", name, false)
--[[ articles.id = 1, cache = false ]]--
```

The result is always a table pointer. You can use cache and still edit the table pointer

# query

Executes a query returning all the results found as a table

```lua
local name = "test"
local articles = db:query("SELECT id FROM articles WHERE name = ?", name, false)
--[[ articles[1].id = 1 ]]--
```

Same rules from `singleQuery` apply to `query`

# execute

Runs the given SQL command. If the command is of the type `INSERT` the last inserted ID is returned

```lua
local name = "test"
local articles = db:execute("UPDATE articles SET name = ? WHERE id = 1", name)
--[[ articles = nil ]]--
```

```lua
local name = "test"
local id = db:execute("INSERT INTO articles (name) VALUES (?)", name)
--[[ id = 1 ]]--