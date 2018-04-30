---
name: Logging
---

# Logging

Castro creates a log file at startup with the following format `year-Month-day.txt`. **A new file is created each day automatically.**

The contents of this file are in text. Each message follows the same structure:

`[LEVEL] (DATE) MESSAGE`

Below is an example on how these messages look like.

```json
[info] (2017-10-29 01:19:29) Encoded map is outdated. Generating new map data 
[info] (2017-10-29 01:19:34) New map data saved to database 
[info] (2017-10-29 01:19:34) House list loaded 
```
