---
Name: file
---

# File metatable

Provides access to file helper functions

- [file:mod(filepath)](#mod)
- [file:exists(filepath)](#exists)
- [file:getFiles(fullpath)](#getfiles)
- [file:getDirectories(fullpath)](#getdirectories)
- [file:unzip(archive, fullpath)](#unzip)

# mod

Gets the last modified time in seconds from the given file.

```lua
local last = file:mod("test.json")
-- last = 1487974045
```

# exists

Checks if the given file exists.

```lua
local exists = file:exists("weird_file.weird")
-- exists = false
```

# getFiles

Returns a table with all the files inside a directory.

```lua
local files = file:getFiles("/home/test")
```

# getDirectories

Returns a table with all the directories inside a directory.

```lua
local files = file:getDirectories("/home/test")
```

# unzip

Extracts a zip archive into the given directory.
Returns true on success and nil, error on failure.

```lua
ok, err = file:unzip("archive.zip", "/home/test")
-- ok == true
-- err == nil
```
