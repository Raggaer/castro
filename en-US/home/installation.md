---
name: Installation
---

# Setup

Castro does not need a HTTP server, this means you dont need Apache, XAMPP, NGINX or anything like that. When you download Castro save it on any folder you wish and then start the application (.exe if using win) this will start the installation process.

> **You need a running MySQL server. Castro does not need a HTTP server**

> **Make sure you run Castro from the main directory (Castro folder).**

# Installation

Castro comes with a very handy installation wizard. If its your first time using Castro start the application, Castro will run on port `:8080` the installation wizard. Just head over to that page and follow the simple steps.

During this process Castro executes your `engine/install.lua` file to set some custom values.

## Conclusion

To sum the whole process here is what you should do in order to start using Castro:

- Get the executable (for your platform) and the needed folders.
- Configure the `config.toml` file
- Start the executable (from the castro directory) and head over to the installation interface (**port: 8080**).
- Follow the installation process (this will encode your map and create a configuration file).
- Restart Castro.