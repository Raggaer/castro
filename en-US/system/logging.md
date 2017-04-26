---
name: Logging
---

# Logging

Castro creates a log file at startup with the following format `year-Month-day.json`. A new file is created each day automatically.

The contents of this file are in JSON. Each message contains two fields:

- level: the warning level.
- message: the logged message.
- time: when the message was logged.

Below is an example on how these messages look like.

```json
{"level":"info","msg":"query: SELECT name, level FROM players ORDER BY level DESC LIMIT 5","time":"2017-04-04T01:48:56+02:00"}
{"level":"info","msg":"query: SELECT COUNT(*) AS count FROM players_online","time":"2017-04-04T01:49:39+02:00"}
{"level":"info","msg":"execute: INSERT INTO castro_onlinechart (count, time) VALUES (0, 1491263379)","time":"2017-04-04T01:49:39+02:00"}
{"level":"info","msg":"execute: DELETE FROM castro_onlinechart WHERE time \u003c 1491262839","time":"2017-04-04T01:49:39+02:00"}
{"level":"info","msg":"query: SELECT COUNT(*) AS count FROM players_online","time":"2017-04-04T01:50:39+02:00"}
```
