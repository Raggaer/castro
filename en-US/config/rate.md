---
name: Rate limit
---

Provides access to the website request rate-limiter.

- [Number](#number)
- [Time](#time)

# Number

The number of requests to trigger the rate-limiter.

# Time

Time in seconds that is allowed between each requests. If the time is less and the number of requests is reached the client will be blocked by the rate-limiter.