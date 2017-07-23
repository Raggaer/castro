---
name: Cookies
---

# Cookies

Provides access to the website cookies options

- [HashKey](#hashkey)"12345678912345671234567891234567"
- [BlockKey](#blockkey)"1234567891234567"
- [Name](#name) = "castro"
- [MaxAge](#maxage) = 1000000

# HashKey

Hash key used to secure store cookies. This value is generated at the installation process. The key must have 32 bits.

# BlockKey

Block key used as block cipher for the cookie storage. This value is generated at installation process. The key must be 16 bits long.

# Name

The current name of the cookie. This value is generated at the installation process with the given format:

```
castro-%s
```

# MaxAge

The cookie max age value