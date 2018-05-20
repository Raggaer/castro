---
name: Static
---

# Static

Provides access to static assets serving options.

- [Enabled](#enabled)
- [Directory](#directory)

# Enabled

Turns the asset server on or off. Usually you turn it off when using an external HTTP server like Caddy or nginx, in this case you dont want Castro to handle the static assets part.

# Directory

The directory where the static content are located, ending with a trailing slash, for example, `public/`