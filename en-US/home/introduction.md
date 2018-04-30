---
name: Introduction
---

# Introduction

Castro is a high performance Open Tibia content management system written in **Go** using **Lua** scripting.

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool.

## Features

Castro provides lua bindings using a pool of lua states. Each request gets a state from the pool. If there are no states available a new one is created and later saved on the pool.

- Standalone application.
    - Castro can run by itself without Apache or NGINX.

- Extensible and solid lua support.
    - Create your own logic using exclusively lua 
    - Everything you need is exposed by Castro.

- Simple installation. 
    - One click process. 
    - Everything is gathered from your config and map files.   
    - Converts normal and znote accounts.

- Plugin manager. 
    - Manage your extensions through Castro itself.

- Clean templates. 
    - Logic is separated from templates. 
    - No more messy files.

- Security
    - All OWASP headers are covered.
    - Highly customizable Content Security Policy.
    - Prepared statements.
    - Template data is escaped by default.

- Shop
    - Highly customizable web integrated shop.
    - PayPal payment processing using REST API.
    - Fortumo support.
    - PayGol support.

## Development status

Castro is production ready