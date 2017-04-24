---
name: Behind a proxy
---

# Behind a proxy

The best setup for Castro is to be behind a HTTP server. The recommended servers are:

- [nginx](https://nginx.org/en/)
- [Caddy](https://caddyserver.com/)

Below are some examples on how to setup Castro behind one of these servers.

# nginx

```ini
location / {
    proxy_pass http://localhost:8080
}
```

# Caddy

```ini
localhost:80 {
	tls email@email.com
	proxy / http://localhost:8080
}
```

This will listen on `localhost:80` and Castro should listen on `localhost:8080`. It is recommended to use the  `tls` setting to enable HTTPS.

You must have `SSL.Proxy = true` on your `config.toml`. For more information about the proxy directive [head to the Caddy docs](https://caddyserver.com/docs/proxy)