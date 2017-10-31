---
name: Rate limit
---

# Rate-limiter

Provides access to the website request rate-limiter.

- [Enabled](#enabled)
- [Number](#number)
- [Time](#time)

# Enabled

Disables or enables the rate-limiter module.

# Number

The number of requests to trigger the rate-limiter.

# Time

Time that is allowed between each requests. If the time is less and the number of requests is reached the client will be blocked by the rate-limiter. This is a time string that follows the [go-duration](https://castroaac.org/docs/config/duration) format.